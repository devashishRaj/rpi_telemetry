-- Create schema
CREATE SCHEMA telemetry;

-- Create the devices table
CREATE TABLE telemetry.devices (
    MacAddress VARCHAR(25) PRIMARY KEY NOT NULL,
    PrivateIP VARCHAR(25) NOT NULL,
    PublicIP VARCHAR(25) NOT NULL,
    Hostname VARCHAR(25) NOT NULL,
    OSType VARCHAR(25) NOT NULL,
    TotalMemory INT NOT NULL
);

-- Create the rpi4b_metrics table
CREATE TABLE telemetry.rpi4b_metrics (
    MacAddress VARCHAR(25) NOT NULL,
    CPUuserLoad DOUBLE PRECISION NOT NULL,
    MemoryUsage INT NOT NULL,
    Temperature REAL NOT NULL,
    TotalProcesses INT NOT NULL,
    TimeStamp TIMESTAMP UNIQUE NOT NULL,
    CONSTRAINT fk_MacAddress FOREIGN KEY (MacAddress) REFERENCES telemetry.devices(MacAddress)
);

-- Create the rpi_temp_alert table
CREATE TABLE telemetry.rpi_temp_alert (
    MacAddress VARCHAR(25) NOT NULL,
    CPUuserLoad DOUBLE PRECISION NOT NULL,
    MemoryUsage INT NOT NULL,
    PrivateIP VARCHAR(25) NOT NULL,
    Temperature REAL NOT NULL,
    TotalProcesses INT NOT NULL,
    TimeStamp TIMESTAMP NOT NULL,
    CONSTRAINT fk_MacAddress FOREIGN KEY (MacAddress) REFERENCES telemetry.devices(MacAddress)
);

