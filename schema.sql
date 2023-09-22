-- Create schema
create schema telemetry;
-- Create the devices table
CREATE TABLE telemetry.devices (
    MacAddress VARCHAR(25) PRIMARY KEY,
    PrivateIP VARCHAR(25),
    PublicIP VARCHAR(25),
    Hostname VARCHAR(25),
    OSType VARCHAR(25),
    TotalMemory INT
);

-- Create the rpi4b_metrics table
CREATE TABLE telemetry.rpi4b_metrics (
    MacAddress VARCHAR(25),
    CPUuserLoad DOUBLE PRECISION,
    MemoryUsage INT,
    Temperature REAL,
    TotalProcesses INT,
    TimeStamp TIMESTAMP UNIQUE,
    CONSTRAINT fk_MacAddress FOREIGN KEY (MacAddress) REFERENCES telemetry.devices(MacAddress)
);

-- Create the rpi_temp_alert table
CREATE TABLE telemetry.rpi_temp_alert (
    MacAddress VARCHAR(25),
    CPUuserLoad DOUBLE PRECISION,
    MemoryUsage INT,
    PrivateIP VARCHAR(25),
    Temperature REAL,
    TotalProcesses INT,
    TimeStamp TIMESTAMP,
    CONSTRAINT fk_MacAddress FOREIGN KEY (MacAddress) REFERENCES telemetry.devices(MacAddress)
);

