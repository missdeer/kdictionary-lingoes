/**************************************************************************
**   Author: Fan Yang
**   Email: missdeer@gmail.com
**   License: see the license.txt file
**************************************************************************/
#include "sqlitewriter.h"
#include <QtCore>
#include <QSqlDatabase>
#include <QSqlQuery>
#include <QSqlError>

SqliteWriter::SqliteWriter(const QString &outputFilePath)
    :dbFilePath(outputFilePath)
{
    QSqlDatabase db = QSqlDatabase::database(outputFilePath, true);
    if (!db.isValid()) {
        db = QSqlDatabase::addDatabase("QSQLITE", outputFilePath);
        db.setDatabaseName(outputFilePath);
    }

    if (!db.isOpen())
        db.open();

    if (db.isOpen())
    {
        query = new QSqlQuery(db);
        query->exec("CREATE TABLE dict(id INTEGER PRIMARY KEY AUTOINCREMENT, word TEXT, content TEXT);");

        query->exec("PRAGMA synchronous = OFF");
        query->exec("PRAGMA journal_mode = MEMORY");
    }
}

SqliteWriter::~SqliteWriter()
{
    delete query;
    QSqlDatabase db = QSqlDatabase::database(dbFilePath, false);
    if (db.isOpen())
        db.close();
}

void SqliteWriter::append(const QString &word, const QString &content)
{
    query->prepare("INSERT INTO dict (word, content) VALUES (:word, :content);");

    query->bindValue(":word", word);
    query->bindValue(":content", content);
    if (!query->exec()) {
        qDebug() << query->lastError();
    }
}

void SqliteWriter::start()
{
    QSqlDatabase db = QSqlDatabase::database(dbFilePath, false);
    if (db.isOpen())
        db.transaction();
}

void SqliteWriter::end()
{
    QSqlDatabase db = QSqlDatabase::database(dbFilePath, false);
    if (db.isOpen())
        db.commit();
}
