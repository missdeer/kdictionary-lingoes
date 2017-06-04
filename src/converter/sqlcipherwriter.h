/**************************************************************************
**   Author: Fan Yang
**   Email: missdeer@gmail.com
**   License: see the license.txt file
**************************************************************************/
#ifndef SQLCIPHERWRITER_H
#define SQLCIPHERWRITER_H

#include <QString>

struct sqlite3;

class SqlcipherWriter
{
public:
    explicit SqlcipherWriter(const QString& outputFilePath, const QString& cipherName, const QString& key);
    ~SqlcipherWriter();
    void append(const QString& word, const QString& content);
    void start();
    void end();
private:
    sqlite3 * db_;
    void execDML(QString statement);
};

#endif // SQLCIPHERWRITER_H
