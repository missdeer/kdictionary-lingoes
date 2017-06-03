/**************************************************************************
**   Author: Fan Yang
**   Email: missdeer@gmail.com
**   License: see the license.txt file
**************************************************************************/
#ifndef SQLITEWRITER_H
#define SQLITEWRITER_H

#include <QString>

class SqliteWriter
{
public:
    explicit SqliteWriter(const QString& outputFilePath);
    ~SqliteWriter();
    void append(const QString& word, const QString& content);
private:
};

#endif // SQLITEWRITER_H
