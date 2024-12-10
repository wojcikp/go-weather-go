## GO-WEATHER-GO APP

PL
Aplikacja go-weather-go powstała w związku z moją chęcią do śledzenia danych pogodowych, co weszło mi ostatnio w nawyk. 
Jest napisana w **Go** i podzielona na dwa moduły (**weather-feed** oraz **weather-app**), które komunikują się ze sobą przy pomocy brokera wiadomości (**RabbitMQ**), gdzie na kolejkę trafiają dane pogodowe ze 172 miast Polski. Dzieje się to przy pomocy modułu **weather-feed**. Moduł pozyskuje dane z zewnętrznego API (open-meteo.com). Następnie dane z kolejki są zczytywane przez moduł **weather-hub** i zasilają bazę danych (**ClickHouse**). W momencie gdy wszystkie dane z kolejki trafią do bazy, uruchomione zostaje przeliczanie danych pogodowych w określonych dla nich **oknach czasowych** Są to:
 - miasto z najwyższą średnią temperaturą w ostatnim tygodniu
 - najbardziej deszczowe miasto w ostatnim tygodniu
 - najbardziej deszczowe miasto w ostatnim miesiącu
 - średnia temperatura w Warszawie w ostatnich 2 tygodniach

Po przeliczeniu wyniki są dostępne w formacie JSON pod adresem: **http://localhost:8081/scores**

### Jak uruchomić aplikację
Musisz posiadać zainstalowanego lokalnie Gita i Dockera. 
1. Otwórz okno terminala i pobierz to repozytorium do dowolnej lokalizacji: 
`git clone https://github.com/wojcikp/go-weather-go.git`
2. Przejdź do folderu z repozytorium: `cd go-weather-go`
3. Uruchom aplikację przy pomocy docker-compose: `docker-compose up`
*Zwróć uwagę, porty **5672, 15672, 9000 oraz 8081** nie mogą być zajęte na Twojej maszynie.*

Gdy wszystkie elementy składowe się uruchomią, aplikacja uaktualni dane pogodowe w bazie i na nowo przeliczy wyniki. Dane na potrzeby demonstracji będą uaktualniane bardzo często, co 3 minuty. Przeliczone wyniki powinieneś móc zobaczyć pod adresem: **http://localhost:8081/scores**
<br>

EN
The **go-weather-go** app was created due to my recent habit of tracking weather forecasts and historical weather data. It is written in the **Go** programming language and consists of two main modules: **weather-feed** and **weather-app**, which communicate via the **RabbitMQ** message broker.

The **weather-feed** module gathers weather data from an external API (open-meteo.com) for 172 cities in Poland and sends it to a **RabbitMQ queue**. The **weather-app** module then collects this data from the queue and inserts it into a **ClickHouse database**. Once all data from the queue has been loaded into the database, the module recalculates **weather scores** for specific time periods. These scores are preset and updated regularly. They include:
- city with the highest average temperature in last week
- the most rainy city in last week
- the most rainy city in last month
- average temperature in Warsaw in last 2 weeks

After calculations are complete, the results are available in JSON format at: http://localhost:8081/scores.
### How to Run the App
To run the app, ensure you have Git and Docker installed on your computer.
1. Open a terminal and clone the repository to a directory of your choice:
`git clone https://github.com/wojcikp/go-weather-go.git`
2. Navigate to the newly created directory: `cd go-weather-go`
3. Start the app using Docker Compose: `docker-compose up`
*Note: Ensure that ports **5672, 15672, 9000, and 8081** are available on your machine.*

Once all components are running, the app will update the weather data in the database and recalculate the weather scores. For demonstration purposes, the weather data is updated frequently (every 3 minutes). The current weather scores can be accessed at: http://localhost:8081/scores.