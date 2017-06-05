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
    if(sqlite3_open_v2(outputFilePath.toUtf8(), &db_, SQLITE_OPEN_READWRITE | SQLITE_OPEN_CREATE, NULL) != SQLITE_OK)
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
    if(sqlite3_prepare_v2(db_, "INSERT INTO dict (word, content) VALUES (?, ?);", -1, &stmt, 0) != SQLITE_OK)
    {
        qWarning() << QString::fromUtf8((const char*)sqlite3_errmsg(db_));
        return;
    }

    QByteArray w = word.toUtf8();
    if(sqlite3_bind_text(stmt, 1, w.constData(), w.length(), SQLITE_STATIC))
    {
        qWarning() << QString::fromUtf8((const char*)sqlite3_errmsg(db_));
        return;
    }

    QByteArray c = content.toUtf8();
    if(sqlite3_bind_text(stmt, 2, c.constData(), c.length(), SQLITE_STATIC))
    {
        qWarning() << QString::fromUtf8((const char*)sqlite3_errmsg(db_));
        return;
    }

    if(sqlite3_step(stmt) != SQLITE_DONE)
    {
        qWarning() << QString::fromUtf8((const char*)sqlite3_errmsg(db_));
        return;
    }
    if(sqlite3_finalize(stmt) != SQLITE_OK)
    {
        qWarning() << QString::fromUtf8((const char*)sqlite3_errmsg(db_));
        return;
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
    int nRet;

    do {
        sqlite3_stmt* pVM = nullptr;
        int res = sqlite3_prepare(db_, statement.toUtf8(), -1, &pVM, nullptr);
        if (res != SQLITE_OK)
        {
            qWarning() << QString::fromUtf8((const char*)sqlite3_errmsg(db_));
            break;
        }

        nRet = sqlite3_step(pVM);

        if (nRet == SQLITE_ERROR) {
            qWarning() << QString::fromUtf8((const char*)sqlite3_errmsg(db_));
            sqlite3_finalize(pVM);
            break;
        }
        nRet = sqlite3_finalize(pVM);
    } while (nRet == SQLITE_SCHEMA);
}
