-- Flight Table
CREATE TABLE Flight (
    FlightCode NVARCHAR(10) PRIMARY KEY NOT NULL,
    Departure NVARCHAR(100) NOT NULL,
    Arrival NVARCHAR(100) NOT NULL,
    DurationInMinutes INT NOT NULL,
    DepartureTime DATETIME NOT NULL,
    DepartureDays NVARCHAR(MAX) NOT NULL,
    BasePrice FLOAT NOT NULL,
    CreatedAt DATETIME NOT NULL
)