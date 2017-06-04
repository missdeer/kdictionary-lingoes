/**************************************************************************
**   Author: Fan Yang
**   Email: missdeer@gmail.com
**   License: see the license.txt file
**************************************************************************/
#include "sqlcipherwriter.h"
#include "sqlite3.h"
#include <QtCore>


SqlcipherWriter::SqlcipherWriter(const QString &outputFilePath, const QString& cipherName, const QString& key)
{
    if(sqlite3_open_v2(outputFilePath.toUtf8(), &db_, SQLITE_OPEN_READWRITE, nullptr) != SQLITE_OK)
    {
        qCritical() << QString::fromUtf8((const char*)sqlite3_errmsg(db_));
        return ;
    }

    sqlite3_key(db_, key.toUtf8(), key.length());
    QString statement = QString("PRAGMA cipher = '%1';").arg(cipherName);

    execDML(statement);
    // create table
    execDML("CREATE TABLE dict(id INTEGER PRIMARY KEY AUTOINCREMENT, word TEXT, content TEXT);");

    execDML("PRAGMA synchronous = OFF");
    execDML("PRAGMA journal_mode = MEMORY");
}

SqlcipherWriter::~SqlcipherWriter()
{
    sqlite3_close_v2(db_);
}

void SqlcipherWriter::append(const QString &word, const QString &content)
{
    sqlite3_stmt* stmt;
    int success = 1;
    if(sqlite3_prepare_v2(db_, "INSERT INTO dict (word, content) VALUES (:word, :content);", -1, &stmt, 0) != SQLITE_OK)
    {
        success = 0;
    }
    if(sqlite3_bind_text(stmt, 1, word.toUtf8(), word.length(), SQLITE_STATIC))
        success = -1;

    if(sqlite3_bind_text(stmt, 2, content.toUtf8(), content.length(), SQLITE_STATIC))
        success = -1;

    if(success == 1 && sqlite3_step(stmt) != SQLITE_DONE)
        success = -1;
    if(success != 0 && sqlite3_finalize(stmt) != SQLITE_OK)
        success = -1;

    if(success != 1)
    {
        qWarning() << QString::fromUtf8((const char*)sqlite3_errmsg(db_));
    }
}

void SqlcipherWriter::start()
{
    execDML("BEGIN TRANSACTION;");
}

void SqlcipherWriter::end()
{
    execDML("COMMIT;");
}

void SqlcipherWriter::execDML(QString statement)
{
    char *errmsg;
    if (sqlite3_exec(db_, statement.toUtf8(), NULL, NULL, &errmsg) != SQLITE_OK)
    {
        qCritical() << QString::fromUtf8(errmsg);
        sqlite3_free(errmsg);
    }
}
