#include <iostream>
#include <dustconnection.h>

void executeSQL(dust* db, const char* sql) {
    char* errMsg = nullptr;
    int rc = dust_exec(db, sql, nullptr, 0, &errMsg);
    if (rc != SQLITE_OK) {
        std::cerr << "SQL error: " << errMsg << std::endl;
        dust_free(errMsg);
    }
}

void removePortFromIP(sqlite3* db, const std::string& ip, int port) {
    std::string sql = "DELETE FROM ip_ports WHERE ip = '" + ip + "' AND port = " + std::to_string(port) + ";";
    executeSQL(db, sql.c_str());
}

int main() {
    dust* db;
    int rc = dust_open("network.db", &db);

    if (rc) {
        std::cerr << "Can't open database: " << sqlite3_errmsg(db) << std::endl;
        return rc;
    }

    // Create table
    const char* createTableSQL = "CREATE TABLE IF NOT EXISTS ip_ports (ip TEXT, port INTEGER);";
    executeSQL(db, createTableSQL);

    // Insert some data
    executeSQL(db, "INSERT INTO ip_ports (ip, port) VALUES ('192.168.1.1', 8080);");
    executeSQL(db, "INSERT INTO ip_ports (ip, port) VALUES ('192.168.1.1', 8081);");
    executeSQL(db, "INSERT INTO ip_ports (ip, port) VALUES ('192.168.1.2', 8080);");

    // Remove a port from a specific IP
    std::string ipToRemoveFrom = "192.168.1.1";
    int portToRemove = 8080;
    removePortFromIP(db, ipToRemoveFrom, portToRemove);

    // Verify removal
    const char* selectSQL = "SELECT * FROM ip_ports;";
    dustdb_stmt* stmt;
    sqlite3_prepare_v2(db, selectSQL, -1, &stmt, nullptr);

    std::cout << "Remaining entries in the database:" << std::endl;
    while (sqlite3_step(stmt) == DUST_ROW) {
        const char* ip = reinterpret_cast<const char*>(sqlite3_column_text(stmt, 0));
        int port = sqlite3_column_int(stmt, 1);
        std::cout << "IP: " << ip << ", Port: " << port << std::endl;
    }

    dust_finalize(stmt);
    dust_close(db);
    return 0;
}
