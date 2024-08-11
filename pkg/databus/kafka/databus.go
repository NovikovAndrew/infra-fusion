package kafka

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
)

type Offset int

const (
	OffsetNewest Offset = -1
	OffsetOldest Offset = -2
)

type Config struct {
	Brokers  []string
	Login    string
	Password string
	Offset   Offset
	ClientID string
	Commit   bool
}

type Databus struct {
	logger       *slog.Logger
	saramaConfig *sarama.Config
	client       sarama.Client
}

func New(cfg Config) (*Databus, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid databus config: %w", err)
	}

	saramaConfig := makeSaramaConfig(cfg)
	client, err := sarama.NewClient(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, errors.New("could not connect to provided brokers")
	}

	return &Databus{
		saramaConfig: saramaConfig,
		client:       client,
	}, nil
}

func validateConfig(cfg Config) error {
	if len(cfg.Brokers) == 0 {
		return errors.New("broker hosts should be provided")
	}

	if cfg.Offset != OffsetNewest && cfg.Offset != OffsetOldest {
		return errors.New("initial consumer offset should be either -1 or -2")
	}
	if cfg.ClientID == "" {
		return errors.New("Client ID is empty")
	}

	return nil
}

func makeSaramaConfig(cfg Config) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	// ==================== Общяя конфигурация  ====================
	// Наименование клиента. Заполнять по наименованию сервиса
	saramaConfig.ClientID = cfg.ClientID
	// TODO: пока неизвестно какая версия кафки используется
	//saramaConfig.Version = sarama.V3_6_0_0

	// Размер буфера для фходящего и исходящего канала сообщений
	// Producer и consumer продолжает обработку некоторых сообщений
	// в фоновом режиме, пока работает пользовательский код, что значительно повышает пропускную способность(Defaults to 256.)
	saramaConfig.ChannelBufferSize = 256

	// ==================== Net Сетевое соединение для клиента ====================
	if cfg.Login != "" && cfg.Password != "" {
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.Net.SASL.User = cfg.Login
		saramaConfig.Net.SASL.Password = cfg.Password
	}
	// ==================== Metadata ====================
	// Количество ретраев метаданных.
	// При недоступности брокера или при ребалансировке группы клиент будет
	// запрашивать метаданные.
	// Запрос метаданных происходит при подключении к брокеру и далее каждый раз,
	// когда меняется схема кластера, происходит ребалансировка или случаются
	// сетевые проблемы. Имеет смысл повысить это число.
	// Метаданные - это актуальные адреса брокеров, список топиков, партиции в нем.
	// По умолчанию значение равно 3. Это означает, что клиент будет пытаться
	// три раза получить метаданные с интервалом в 250 миллисекунд и после этого
	// сдастся. Результатом будет сообщение вида `Your metadata is out of date`
	// Идея в том, чтобы увеличить количество попыток получения метаданных, чтобы
	// увеличить шансы самостоятельного восстановления клиентов, т.к. возможны
	// сетевые проблемы, которые, в свою очередь, редко длятся долго, но достаточно,
	// чтобы клиент после третьей попытки завершился с ошибкой и перестал читать
	// сообщения. Если один, несколько брокеров или весь кластер становятся
	// недоступными, клиент будет стараться перезапросить метаданные для того,
	// чтобы продолжить чтение сообщений.
	// Следует согласовывать с настройкой Metadata.Retry.Backoff
	// или Metadata.Retry.BackoffFunc.
	saramaConfig.Metadata.Retry.Max = 100
	// Просим клиент не запрашивать сразу все топики, а только те, которые
	// необходимы клиенту в процессе работы. Эта настройка помогает избежать
	// бесполезного расходования памяти в виде хранения всех топиков, о которых
	// знает брокер.
	// Запрашивать полный список топиков или нет. Если указано false, то клиент
	// будет запрашивать только те топики, которые необходимо. Это экономит память,
	// используемую приложением, т. к. приложение обычно использует всего несколько
	// топиков, а у брокеров может быть их огромное количество. Соответственно, при
	// значении true каждый инстанс клиента будет выкачивать и хранить все метаданные
	// по всем имеющимся у брокеров топикам. Поэтому переопределяем параметр(По умолчанию в sarama, заполняется как true)
	saramaConfig.Metadata.Full = false
	// При необходимости указываем насколько часто клиент должен обновлять
	// метаданные (информацию о брокерах, топиках, разделах, лидерстве).
	saramaConfig.Metadata.RefreshFrequency = 10 * time.Minute
	// Автосоздание топиков
	saramaConfig.Metadata.AllowAutoTopicCreation = false

	// ==================== Producer ====================
	// Дожидаемся подтверждений от всех in-sync реплик для того, чтобы гарантированно
	// не потерять данные после отказа одного из брокеров.
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	// ВНИМАНИЕ! Следует учесть что для асинхронного producer-а эти параметры нужно отключить
	// Либо читать из каналов успешнов отправлений\ошибок. иначе deadlock
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	// Ограничения по размерам сообщений и батчей что отправляет producer
	//config.Producer.Flush.Frequency = 5 * time.Second
	//config.Producer.Flush.Bytes = int(sarama.MaxRequestSize)
	//config.Producer.Flush.Messages = 500

	// ==================== Consumer ====================
	saramaConfig.Consumer.Offsets.Initial = int64(cfg.Offset)

	// Таймаут используется для обнаружения сбоев. Консьюмер периодически
	// отправляет брокеру хертбиты, чтобы сообщить ему о своей активности.
	// Если по истечению этого таймаута брокер не получает хертбит - он удалит
	// консьюмера из группы и начнет ребалансировку.
	// В документации драйвера указано, что это значение должно находиться в
	// диапазоне от `group.min.session.timeout.ms` до `group.max.session.timeout.ms`
	// параметров, указанных в конфигурации брокера.
	saramaConfig.Consumer.Group.Session.Timeout = 10 * time.Second

	// Время между посылкой хертбитов брокеру. Хертбиты используются для того,
	// чтобы сессия консьюмера оставалась активной и брокер не начинал
	// перебалансировку группы в связи с потерей консьюмера.
	// Значение должно быть меньше значения Consumer.Group.Session.Timeout
	// и обычно не должно превышать 1/3 от него.
	saramaConfig.Consumer.Group.Heartbeat.Interval = 3 * time.Second

	// Максимальное время, за которое сам консьюмер ожидает, что сообщение
	// обработано бизнес-логикой приложения. Если запись в канал принимаемых
	// консьюмером сообщений занимает больше времени, чем указано в параметре,
	// консьюмер-группа перестанет получать сообщения из партиции до тех пор,
	// пока место в нем не освободится.
	// Стоит обратить внимание, что т. к. канал сообщений буферизированный,
	// фактически время отсрочки составляет (MaxProcessingTime * ChannelBufferSize).
	// saramaConfig.Consumer.MaxProcessingTime = 100 * time.Millisecond

	// Когда новый консьюмер присоединяется к группе, группа пытается
	// перебалансировать нагрузку, чтобы назначить разделы каждому консьюмеру.
	// Если группа изменится во время выполнения этого назначения, повторная
	// ребалансировка завершится ошибкой и повторится. Этот параметр определяет
	// максимальное количество попыток до отказа (по умолчанию 4).
	saramaConfig.Consumer.Group.Rebalance.Retry.Max = 4
	saramaConfig.Consumer.Group.Rebalance.Retry.Backoff = 2 * time.Second
	saramaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	// Время между попытками консьюмера присоединиться к группе в случае неудачи.
	// conf.Consumer.Group.Rebalance.Retry.Backoff = 2 * time.Second

	// Время между неудачными попытками считать сообщения из партиции.
	//saramaConfig.Consumer.Retry.Backoff = 2 * time.Second
	//saramaConfig.Consumer.Fetch.Min = 100 * 1000 // 100 Kb
	//saramaConfig.Consumer.MaxWaitTime = 500 * time.Millisecond // default 250 ms

	// Ошибки что происходят в consumer-е возвращаются в отдельный канал.
	// вид ошибок - проблема с брокерами.
	// заблокированы партиции\брокеры
	// ErrOffsetMetadataTooLarge, ErrInvalidCommitOffsetSize, ErrFencedInstancedId другие
	// по дефолту, не критичные ошибки так же будут попадать
	saramaConfig.Consumer.Return.Errors = true

	return saramaConfig
}

func (d *Databus) Close() error {
	return d.client.Close()
}
