const mongoHost = process.env.INCIDENT_LOG_API_MONGODB_HOST;
const mongoPort = process.env.INCIDENT_LOG_API_MONGODB_PORT;

const mongoUser = process.env.INCIDENT_LOG_API_MONGODB_USERNAME;
const mongoPassword = process.env.INCIDENT_LOG_API_MONGODB_PASSWORD;

const database = process.env.INCIDENT_LOG_API_MONGODB_DATABASE;
const collection = process.env.INCIDENT_LOG_API_MONGODB_COLLECTION;

const retrySeconds = parseInt(process.env.RETRY_CONNECTION_SECONDS || "5") || 5;

let connection;

while (true) {
    try {
        const auth =
            mongoUser && mongoPassword
                ? `${mongoUser}:${mongoPassword}@`
                : "";

        connection = Mongo(`mongodb://${auth}${mongoHost}:${mongoPort}`);
        break;
    } catch (exception) {
        print(`Cannot connect to mongoDB: ${exception}`);
        print(`Will retry after ${retrySeconds} seconds`);
        sleep(retrySeconds * 1000);
    }
}

const databases = connection.getDBNames();

if (databases.includes(database)) {
    const dbInstance = connection.getDB(database);
    const collections = dbInstance.getCollectionNames();

    if (collections.includes(collection)) {
        print(`Collection '${collection}' already exists in database '${database}'`);
        process.exit(0);
    }
}

const db = connection.getDB(database);

db.createCollection(collection);

db[collection].createIndex({ id: 1 });

let result = db[collection].insertMany([
    {
        id: "hospital-security-log",
        name: "Nemocničný bezpečnostný denník",
        location: "Univerzitná nemocnica",
        incidents: [
            {
                id: "INC-001",
                incidentType: "Bezpečnostná udalosť",
                location: "Urgentný príjem",
                occurredAt: new Date("2038-12-24T10:05:00.000Z"),
                description: "Neoprávnený vstup do vyhradenej zóny.",
                severity: "Vysoká",
                status: "Nový",
                attachments: ["kamera-zaznam.mp4"],
                investigationReport: "",
                notes: "Incident čaká na preverenie.",
            },
        ],
        predefinedIncidentTypes: [
            {
                value: "Bezpečnostná udalosť",
                code: "security-breach",
                typicalSeverity: "Vysoká",
                description: "Neoprávnený vstup alebo porušenie bezpečnostných pravidiel.",
            },
            {
                value: "Technický incident",
                code: "technical-incident",
                typicalSeverity: "Stredná",
                description: "Výpadok alebo porucha technického zariadenia.",
            },
            {
                value: "Pád pacienta",
                code: "patient-fall",
                typicalSeverity: "Stredná",
                description: "Pád pacienta v priestoroch nemocnice.",
            },
            {
                value: "Chyba pri liečbe",
                code: "medication-error",
                typicalSeverity: "Kritická",
                description: "Nesprávne podanie alebo evidencia lieku.",
            },
        ],
    },
]);

if (result.writeError) {
    console.error(result);
    print(`Error when writing the data: ${result.errmsg}`);
}

process.exit(0);