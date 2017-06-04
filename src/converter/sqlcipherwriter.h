/**************************************************************************
**   Author: Fan Yang
**   Email: missdeer@gmail.com
**   License: see the license.txt file
**************************************************************************/
#ifndef SQLCIPHERWRITER_H
#define SQLCIPHERWRITER_H

#include <QString>

class SqlcipherWriter
{
public:
    explicit SqlcipherWriter(const QString& outputFilePath, const QString& cipherName, const QString& key);
    ~SqlcipherWriter();
    void append(const QString& word, const QString& content);
    void start();
    void end();
};

#endif // SQLCIPHERWRITER_H
