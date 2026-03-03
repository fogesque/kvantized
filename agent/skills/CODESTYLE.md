# C++ Code Style and Development Rules

## CS-01: Парадигма

### Сочетание ООП и ФП

При разработке используется комбинация объектно-ориентированного и функционального подходов. От ООП берётся объектное проектирование и инкапсуляция данных, от ФП — ясность последовательности обработки данных, чистые функции и неизменяемость.

Принципы:
- **Классы** инкапсулируют данные и предоставляют методы для работы с ними
- **Функции обработки данных** должны быть чистыми: принимать данные, возвращать результат, не мутировать входные аргументы
- **Методы трансформации** возвращают новый объект вместо изменения текущего
- **Композиция функций** используется для построения цепочек обработки: фильтрация → трансформация → агрегация

```cpp
// ООП: инкапсуляция данных
class Order
{
public:
    /// @brief Creates order with given id and amount
    static std::tuple<OrderPtr, error> Create(int id, double amount);

    /// @brief Returns order id
    int GetId() const;
    /// @brief Returns order amount
    double GetAmount() const;
    /// @brief Returns payment status
    bool GetIsPaid() const;

    /// @brief Marks order as paid
    void MarkAsPaid();

    /// @brief Returns new order with discounted amount (does not modify current object)
    Order WithDiscount(double percentage) const;

    // ...
};

// ФП: чистые статические функции для обработки коллекций
class OrderProcessor
{
public:
    /// @brief Filters unpaid orders (pure function)
    static std::vector<Order> FilterUnpaid(const std::vector<Order> & orders);

    /// @brief Applies discount to all orders, returns new collection (pure function)
    static std::vector<Order> ApplyDiscount(const std::vector<Order> & orders, double discount);

    /// @brief Calculates total amount (pure function)
    static double CalculateTotal(const std::vector<Order> & orders);

    /// @brief Composes processing pipeline: filter → discount → total
    static double ProcessOrders(const std::vector<Order> & orders);
};

// Использование: ясная последовательность обработки
double total = OrderProcessor::ProcessOrders(orders);
```

### Принцип единой ответственности

Каждый класс и каждая функция должны иметь одну ответственность. Класс отвечает за одну сущность или одну задачу. Функция выполняет одну операцию. Если название класса или функции содержит союз «и» или перечисление действий — это признак нарушения принципа.

```cpp
// Плохо: класс совмещает бизнес-логику и форматирование
class ReportGenerator
{
public:
    std::tuple<Report, error> GenerateReport(const Data & data);
    error SendReportByEmail(const Report & report, const std::string & recipient);
    std::string FormatReportAsHtml(const Report & report);
};

// Хорошо: каждый класс — одна ответственность
class ReportGenerator
{
public:
    std::tuple<Report, error> Generate(const Data & data);
};

class ReportFormatter
{
public:
    static std::string ToHtml(const Report & report);
    static std::string ToCsv(const Report & report);
};

class ReportSender
{
public:
    static error SendByEmail(const Report & report, const std::string & recipient);
};
```

Для функций — аналогично. Функция должна выполнять одно действие. Если требуется последовательность шагов — она должна быть обёрнута в функцию с названием, отражающим единый смысл этой последовательности.

```cpp
// Плохо: название функции перечисляет действия
error ValidateAndSaveConfig(const Config & config);
error ParseAndExecuteCommand(const std::string & input);

// Хорошо: одно действие — одна функция
error ValidateConfig(const Config & config);
error SaveConfig(const Config & config);

// Хорошо: если нужна последовательность — название отражает единую цель
error ApplyConfig(const Config & config);
```

## CS-02: Использование языка C++

Основной тезис — не использовать сложные фичи языка ради них самих. Применять только те возможности C++, которые не портят читаемость и упрощают архитектуру.

### `auto`, `const`, `constexpr`

Использовать `auto` для вывода типов везде, где это возможно. Использовать `const` для всех переменных и параметров, которые не изменяются. Использовать `constexpr` для всех значений, вычислимых на этапе компиляции.

```cpp
// Плохо: явные типы, отсутствие const
std::string name = config.GetName();
std::vector<Order> orders = store::GetOrders();
int maxSize = 1024;

// Хорошо: auto, const, constexpr
const auto name = config.GetName();
const auto orders = store::GetOrders();
constexpr auto maxSize = 1024;
```

### Range-based for

Для итерации по коллекциям всегда использовать range-based `for`. Элементы передаются по `const auto &`, если не требуется модификация.

```cpp
// Плохо: итерация по индексу без необходимости
for (std::size_t i = 0; i < orders.size(); ++i) {
    logging::Log("Order: {}", orders[i].GetId());
}

// Хорошо: range-based for
for (const auto & order : orders) {
    logging::Log("Order: {}", order.GetId());
}
```

### Lambda-выражения

Lambda-выражения допустимы там, где они не затрудняют понимание кода и не нарушают его структуру. Подходят для коротких callback-ов, алгоритмов STL и группировки локальной последовательности операций.

```cpp
// Хорошо: lambda в алгоритме STL
const auto unpaid = std::views::filter(orders, [](const auto & order) {
    return !order.GetIsPaid();
});

// Хорошо: lambda для группировки fallible-операций
const auto runSetup = [&module]() -> error {
    // Initialize module resources
    auto err = device::InitResources(module);
    if (err) {
        return errors::Wrap(err, "Failed to init resources");
    }

    // Apply module configuration
    err = device::ApplyConfig(module);
    if (err) {
        return errors::Wrap(err, "Failed to apply config");
    }

    return nullptr;
};
```

### STL и проверенные библиотеки

Использовать STL и зарекомендовавшие себя библиотеки (vcpkg) вместо написания собственных реализаций. Не изобретать велосипеды для структур данных, алгоритмов, работы с потоками и т.д.

```cpp
// Плохо: собственная реализация поиска
bool found = false;
for (const auto & item : items) {
    if (item.GetId() == targetId) {
        found = true;
        break;
    }
}

// Хорошо: STL-алгоритм
const auto found = std::ranges::any_of(items, [&targetId](const auto & item) {
    return item.GetId() == targetId;
});
```

### Запрет magic numbers

Литеральные числовые значения в коде запрещены. Все числовые константы должны быть именованы через `constexpr` или `inline constexpr` в соответствующем namespace.

```cpp
// Плохо: magic numbers
if (retryCount > 3) {
    return errors::New("Too many retries");
}

std::this_thread::sleep_for(std::chrono::milliseconds(500));

// Хорошо: именованные константы
namespace policy
{
inline constexpr auto MaxRetries = 3;
inline constexpr auto RetryDelay = std::chrono::milliseconds(500);
}  // namespace policy

if (retryCount > policy::MaxRetries) {
    return errors::New("Too many retries");
}

std::this_thread::sleep_for(policy::RetryDelay);
```

### Обязательные фигурные скобки

Однострочные `if`, `else`, `for`, `while` без фигурных скобок запрещены. Тело управляющих конструкций всегда оборачивается в `{}`, даже если содержит одно выражение.

```cpp
// Плохо: без фигурных скобок
if (err)
    return err;

for (const auto & item : items)
    pipeline::Process(item);

if (status.IsReady())
    device::Start(module);
else
    device::Reset(module);

// Хорошо: фигурные скобки всегда
if (err) {
    return err;
}

for (const auto & item : items) {
    pipeline::Process(item);
}

if (status.IsReady()) {
    device::Start(module);
} else {
    device::Reset(module);
}
```

### Запрет исключений

Исключения (`throw`, `try`/`catch`) не используются. Для обработки ошибок применяется подход с возвратом `error` (см. раздел «Обработка ошибок»).

```cpp
// Плохо: исключения
std::tuple<Config, error> ParseConfig(const std::string & path)
{
    try {
        auto data = io::ReadFile(path);
        return {Config{data}, nullptr};
    } catch (const std::exception & exception) {
        return { {}, errors::New(exception.what())};
    }
}

// Хорошо: возврат ошибки
std::tuple<Config, error> ParseConfig(const std::string & path)
{
    // Read config file
    auto [data, err] = io::ReadFile(path);
    if (err) {
        return { {}, errors::Wrap(err, "Failed to read config")};
    }

    return {Config{data}, nullptr};
}
```

## CS-03: Обработка ошибок

### Возврат ошибок из функций

Для обработки ошибок используется библиотека [errors](https://github.com/fogesque/errors). Тип `error` — это алиас для `std::shared_ptr<errors::Error>`. Любая функция, которая может технически завершиться с ошибкой, должна возвращать `error`. Успешное выполнение обозначается возвратом `nullptr`, ошибка — ненулевым указателем.

Если функция возвращает данные вместе с ошибкой, используется `std::tuple` и структурированные привязки. Объект ошибки всегда последний в кортеже.

```cpp
#include <errors/errors.hpp>

error SaveConfig(const Config & config)
{
    // Open config file for writing
    auto [file, err] = io::OpenFile(config.path);
    if (err) {
        return errors::Wrap(err, "Failed to save config");
    }

    // Write config data to file
    err = io::WriteData(file, config.data);
    if (err) {
        return errors::Wrap(err, "Failed to save config");
    }

    return nullptr;
}
```

### Формат сообщений об ошибках

Каждое сообщение об ошибке должно начинаться с заглавной буквы.

```cpp
// Плохо: сообщение с маленькой буквы
return errors::New("connection refused");
return errors::Wrap(err, "failed to connect");

// Хорошо: сообщение с заглавной буквы
return errors::New("Connection refused");
return errors::Wrap(err, "Failed to connect");
return errors::Errorf("Device '{}' not found on port {}", device.Name(), port);
```

### Code flow при обработке ошибок

Каждый вызов fallible-функции отделяется пустой строкой от предыдущего блока. Перед вызовом ставится комментарий, описывающий логический этап. Это снижает визуальную загруженность кода. Используется только ранний возврат (`if (err) { return ... }`) — конструкции `else` не применяются.

```cpp
error InitializeModule(DeviceModulePtr module)
{
    // Setup module with default configuration
    auto err = device::SetupModule(module);
    if (err) {
        return errors::Wrap(err, "Failed to setup module");
    }

    // Start module after successful setup
    err = device::StartModule(module);
    if (err) {
        return errors::Wrap(err, "Failed to start module");
    }

    // Get module status after its start
    auto [status, err2] = device::GetModuleStatus(module);
    if (err2) {
        return errors::Wrap(err2, "Failed to get module status");
    }

    logging::Log("Module status: {}", status.Message());
    return nullptr;
}
```

## CS-04: Проектирование функций

### Чистота функций и возвращаемые значения

Функции должны следовать принципу чистоты: не влиять на внешние объекты и не использовать выходные параметры в аргументах. Все результаты возвращаются через возвращаемое значение.

Если функция возвращает несколько объектов, используется `std::tuple`. Использование `std::pair` запрещено — всегда `std::tuple`, даже для двух элементов. Если функция может завершиться с ошибкой, объект `error` всегда последний элемент кортежа.

```cpp
// Плохо: выходной параметр в аргументах
error ParseConfig(const std::string & path, Config & out_config);

// Плохо: std::pair
std::pair<Config, error> ParseConfig(const std::string & path);

// Хорошо: std::tuple, ошибка последняя
std::tuple<Config, error> ParseConfig(const std::string & path);

// Хорошо: несколько возвращаемых значений
std::tuple<Host, Port, error> ParseAddress(const std::string & address);
```

Вызов таких функций с использованием структурированных привязок:

```cpp
auto [config, err] = parser::ParseConfig("/etc/app/config.json");
if (err) {
    return errors::Wrap(err, "Failed to parse config");
}

auto [host, port, err2] = network::ParseAddress(config.address);
if (err2) {
    return errors::Wrap(err2, "Failed to parse address");
}
```

### Ограничение количества параметров и возвращаемых значений

Список аргументов функций и конструкторов не должен превышать три элемента. Если требуется больше — параметры упаковываются в структуры, сгруппированные по смыслу.

Аналогично для возвращаемых кортежей: максимум два значения (не считая `error`). Если нужно вернуть больше — результат упаковывается в структуру.

```cpp
// Плохо: слишком много аргументов
error CreateConnection(const std::string & host, int port, int timeout,
                       bool useTls, int maxRetries);

// Хорошо: параметры сгруппированы в структуру
struct ConnectionConfig {
    std::string host;
    int port = 0;
    int timeout = 0;
    bool useTls = false;
    int maxRetries = 0;
};

error CreateConnection(const ConnectionConfig & config);

// Плохо: слишком много значений в кортеже (3 значения + error)
std::tuple<Host, Port, Protocol, error> ParseEndpoint(const std::string & endpoint);

// Хорошо: результат упакован в структуру
struct Endpoint {
    Host host;
    Port port;
    Protocol protocol;
};

std::tuple<Endpoint, error> ParseEndpoint(const std::string & endpoint);

// Хорошо: два значения + error допустимо
std::tuple<Host, Port, error> ParseAddress(const std::string & address);
```

### Умные указатели

Объекты классов всегда создаются и передаются через `std::shared_ptr`. Каждый класс должен иметь соответствующий type alias `ClassNamePtr`. Сырые указатели (`new`/`delete`) не используются.

```cpp
// Плохо: сырые указатели
auto module = new DeviceModule();
delete module;

// Хорошо: умные указатели через фабричный метод
auto [module, err] = DeviceModule::Create(config);

// Хорошо: передача через умный указатель
error StartModule(DeviceModulePtr module);
```

## CS-05: Именование

### Именование методов классов

Публичные методы именуются в стиле PascalCase — с заглавной буквы. Приватные методы именуются в стиле camelCase — с маленькой буквы.

```cpp
class ConnectionPool
{
public:
    std::tuple<Connection, error> Acquire();
    error Release(Connection connection);
    std::size_t Size() const;

private:
    error validateConnection(const Connection & connection);
    void removeExpired();
};
```

### Именование интерфейсов

Интерфейсные классы (содержащие только чисто виртуальные методы) именуются с префиксом `I`. Класс может реализовывать несколько интерфейсов.

```cpp
// Интерфейс модуля
class IModule
{
public:
    virtual error Start() = 0;
    virtual void Stop() = 0;
};

// Интерфейс свойства
class IAwaitable
{
public:
    virtual Awaitable Await() = 0;
};

// Реализация нескольких интерфейсов
class DeviceModule : public IModule, public IAwaitable
{
public:
    error Start() override;
    void Stop() override;
    Awaitable Await() override;
};
```

### Именование и доступ к членам класса

Члены класса именуются без суффиксов (`_`) и префиксов (`m_`) — обычными именами. Для обращения к любым членам и методам класса всегда используется `this->`.

```cpp
// Плохо: суффикс, префикс, обращение без this
class Server
{
private:
    std::string m_host;
    int port_;

    void setup()
    {
        m_host = "localhost";
        port_ = 8080;
    }
};

// Хорошо: обычные имена, доступ через this->
class Server
{
public:
    error Start();

private:
    std::string host;
    int port = 0;

    void setup()
    {
        this->host = "localhost";
        this->port = 8080;
        this->Start();
    }
};
```

### Инициализация членов класса и полей структур

Все члены класса и поля структур простых типов (`int`, `bool`, `double`, `enum` и т.д.) должны быть проинициализированы значением по умолчанию при объявлении. Все умные указатели также должны быть явно проинициализированы (`= nullptr`). Сложные типы с собственными конструкторами по умолчанию (`std::string`, `std::vector`, `std::map` и т.д.) инициализировать не нужно.

```cpp
// Плохо: примитивные типы и указатели без инициализации
class Worker
{
private:
    bool running;
    int retryCount;
    DevicePtr device;
    std::string name;
};

// Хорошо: примитивные типы и указатели инициализированы, сложные типы — нет
class Worker
{
private:
    bool running = false;
    int retryCount = 0;
    DevicePtr device = nullptr;
    std::string name;
};
```

### Именование переменных и типов

Названия переменных и типов должны быть читаемыми и понятными. Однобуквенные переменные запрещены, за исключением счётчиков циклов (`i`, `j`, `k`). Слова в названиях пишутся полностью — сокращения допустимы только если они общеприняты и однозначны.

Допустимые сокращения: `config` (configuration), `ptr` (pointer), `err` (error), `arg` / `args` (argument/arguments) и подобные широко известные.

Недопустимые сокращения: `addr` (address), `conn` (connection), `msg` (message), `btn` (button), `srv` (server) и прочие, которые не являются общеупотребительными в индустрии.

```cpp
// Плохо: однобуквенные и неочевидные сокращения
auto c = network::CreateConnection();
auto addr = network::ParseAddress(input);
auto msg = messaging::BuildMessage(data);

// Хорошо: полные и читаемые названия
auto connection = network::CreateConnection();
auto address = network::ParseAddress(input);
auto message = messaging::BuildMessage(data);

// Хорошо: допустимые сокращения
auto [config, err] = parser::ParseConfig(path);
auto ptr = std::make_shared<Handler>();

// Хорошо: счётчик в цикле
for (std::size_t i = 0; i < items.size(); ++i) {
    pipeline::Process(items[i]);
}
```

### Аббревиатуры в именах

Аббревиатуры в именах переменных, типов, функций и методов записываются как обычные слова — с заглавной первой буквой и строчными остальными. Это обеспечивает единообразие стиля PascalCase / camelCase и улучшает читаемость на стыке слов. Правило применяется к аббревиатурам любой длины.

```cpp
// Плохо: аббревиатура полностью в верхнем регистре
class HTTPClient
{
public:
    std::tuple<APIResponse, error> SendHTTPRequest(const std::string & url);

private:
    std::string apiURL;
    int userID = 0;
};

// Хорошо: аббревиатура как обычное слово
class HttpClient
{
public:
    std::tuple<ApiResponse, error> SendHttpRequest(const std::string & url);

private:
    std::string apiUrl;
    int userId = 0;
};
```

### Явная квалификация вызовов

В коде не должно встречаться вызовов функций или обращений к переменным просто по имени, если они не являются локальными. Обращение по голому имени допустимо только для локальных переменных и параметров функции. Это повышает читаемость и позволяет сразу понять, откуда берётся вызываемая функция или переменная.

- Члены и методы класса — через `this->`
- Глобальные и namespace-функции — через явное указание `namespace::Function()`
- Статические методы другого класса — через `ClassName::Method()`

```cpp
// Плохо: непонятно откуда берутся функции
error Initialize(DeviceModulePtr module)
{
    auto err = SetupModule(module);
    if (err) {
        return Wrap(err, "Failed to setup");
    }

    Log("Module initialized");
    return nullptr;
}

// Хорошо: явная квалификация каждого вызова
error Initialize(DeviceModulePtr module)
{
    // Setup device module
    auto err = device::SetupModule(module);
    if (err) {
        return errors::Wrap(err, "Failed to setup");
    }

    logging::Log("Module initialized");
    return nullptr;
}
```

## CS-06: Глобальная область видимости

### Namespace, функции, переменные и константы

Следует избегать создания глобальных функций, переменных и констант. При необходимости все глобальные объекты должны быть размещены в `namespace`. Глобальные переменные объявляются как `inline` для предотвращения ошибок множественного определения при линковке. Все глобальные имена (функции, переменные, константы) начинаются с заглавной буквы.

```cpp
// Плохо: глобальные объекты без namespace
const int maxRetries = 3;
std::string defaultHost = "localhost";
error Connect(const std::string & address);

// Хорошо: всё в namespace, имена с заглавной буквы, переменные inline
namespace network
{

inline constexpr int MaxRetries = 3;
inline std::string DefaultHost = "localhost";

error Connect(const std::string & address);

}  // namespace network
```

## CS-07: Структура файлов

### Именование и организация файлов

Каждый заголовочный файл (`.hpp`) именуется по названию главного класса, который в нём объявлен. Исходный файл (`.cpp`) с реализацией именуется аналогично. Имя файла совпадает с именем класса в PascalCase.

Заголовочные и исходные файлы должны быть разделены по разным директориям.

```
project/
├── include/
│   └── fight/
│       ├── FightClub.hpp
│       └── Fighter.hpp
└── src/
    └── fight/
        ├── FightClub.cpp
        └── Fighter.cpp
```

### Структура заголовочного файла

Заголовочный файл должен быть оформлен в строгой последовательности. Каждый элемент структуры описан ниже.

#### Forward declarations

В начале namespace, до определения класса, размещаются forward declarations для типов, используемых в заголовке. Это уменьшает количество включаемых заголовков и ускоряет компиляцию.

#### Type aliases для указателей

После forward declarations объявляются алиасы для `std::shared_ptr` на класс. Это упрощает использование класса в остальном коде.

```cpp
// Forward declarations
class FightClub;

// Type aliases
using FightClubPtr = std::shared_ptr<FightClub>;
```

#### Константы

Если класс использует именованные константы ошибок или другие константы, они размещаются в отдельном вложенном `namespace` перед определением класса.

```cpp
namespace ErrorTypes
{
inline const auto FighterNotFound = errors::New("Fighter not found");
}  // namespace ErrorTypes
```

#### Фабричный метод

Класс должен предоставлять статический фабричный метод `Create` для создания экземпляров. Фабричный метод всегда располагается первым в `public` секции. Существует две формы:

1. **Fallible** — создание может завершиться ошибкой. Возвращает `std::tuple<Ptr, error>`:

```cpp
static std::tuple<FightClubPtr, error> Create(const Config & config);
```

2. **Infallible** — создание всегда успешно. Возвращает только указатель (никогда не `nullptr`):

```cpp
static FightClubPtr Create(const Config & config);
```

#### Логические секции с комментариями

Методы и переменные внутри класса группируются в логические секции. Каждая секция отделяется комментарием вида `/// [Название секции]`.

#### Комментарии к методам и переменным

Каждый публичный и приватный метод, а также каждая переменная-член должны иметь документирующий комментарий `/// @brief`.

#### `#pragma region` для сворачивания секций

Секции конструкторов/деструкторов и приватных методов оборачиваются в `#pragma region` / `#pragma endregion` для удобства сворачивания в IDE. Имя региона формируется как `ClassName::SectionName`.

#### Полная последовательность элементов

1. `#pragma once`
2. Включение заголовков (`#include`)
3. `namespace`
4. Forward declarations
5. Type aliases
6. Константы (если нужны)
7. Документирующий комментарий класса (`/// @brief`)
8. Определение класса:
   - `public:`
     - `/// [Fabric Methods]` — фабричные методы
     - `/// [Секция]` — логические группы публичных методов
     - `/// [Construction & Destruction]` — конструкторы/деструктор в `#pragma region`
   - `private:`
     - Приватные методы в `#pragma region`, разделённые на логические секции
     - `/// [Properties]` — переменные-члены, разделённые на логические секции

#### Полный пример заголовочного файла

```cpp
// include/fight/FightClub.hpp
#pragma once

#include <errors/errors.hpp>

namespace fight
{

// Forward declarations
class FightClub;

// Type aliases
using FightClubPtr = std::shared_ptr<FightClub>;

/// @brief Errors that can occur in FightClub
namespace ErrorTypes
{
inline const auto FighterNotFound = errors::New("Fighter not found");
}  // namespace ErrorTypes

///
/// @brief
/// FightClub manages fighters and provides operations for adding and searching fighters
///
class FightClub
{
public:
    /// [Fabric Methods]

    /// @brief Creates FightClub instance with given configuration
    static std::tuple<FightClubPtr, error> Create(const std::string & name);

    /// [Fighter Management]

    /// @brief Adds fighter to the club
    error AddFighter(const Fighter & fighter);
    /// @brief Finds fighter by name
    std::tuple<Fighter, error> FindFighter(const std::string & name);

    /// [Construction & Destruction]

#pragma region FightClub::Construct

    /// @brief Copy constructor is deleted
    FightClub(const FightClub &) = delete;
    /// @brief Copy operator is deleted
    FightClub & operator=(const FightClub &) = delete;

    /// @brief Config struct for object construction
    struct Config {
        std::string name;
    };

    /// @brief Constructor
    /// @warning Avoid using this constructor since class has static fabric methods
    explicit FightClub(const Config & config);
    /// @brief Destructor
    ~FightClub();

#pragma endregion

private:
#pragma region FightClub::PrivateMethods

    /// [Validation]

    /// @brief Checks if fighter with given name already exists
    bool fighterExists(const std::string & name);

#pragma endregion

    /// [Properties]

    /// @brief Club name
    std::string name;
    /// @brief List of fighters in the club
    std::vector<Fighter> fighters;
};

}  // namespace fight
```

```cpp
// src/fight/FightClub.cpp
#include <fight/FightClub.hpp>

namespace fight
{

std::tuple<FightClubPtr, error> FightClub::Create(const std::string & name)
{
    const auto config = Config{.name = name};
    const auto club = std::make_shared<FightClub>(config);
    return {club, nullptr};
}

FightClub::FightClub(const Config & config)
{
    this->name = config.name;
}

FightClub::~FightClub() = default;

error FightClub::AddFighter(const Fighter & fighter)
{
    if (this->fighterExists(fighter.Name())) {
        return errors::Errorf("Fighter '{}' already exists", fighter.Name());
    }

    this->fighters.push_back(fighter);
    return nullptr;
}

std::tuple<Fighter, error> FightClub::FindFighter(const std::string & name)
{
    const auto result = std::ranges::find_if(this->fighters, [&name](const auto & fighter) {
        return fighter.Name() == name;
    });

    if (result == this->fighters.end()) {
        return { {}, ErrorTypes::FighterNotFound};
    }

    return {*result, nullptr};
}

bool FightClub::fighterExists(const std::string & name)
{
    return std::ranges::any_of(this->fighters, [&name](const auto & fighter) {
        return fighter.Name() == name;
    });
}

}  // namespace fight
```
