/**************************************************************************
**   Author: Fan Yang
**   Email: missdeer@gmail.com
**   License: see the license.txt file
**************************************************************************/
#ifndef PLAINTEXTWRITER_H
#define PLAINTEXTWRITER_H

#include <QFile>
#include <QTextStream>

class PlainTextWriter
{
public:
    explicit PlainTextWriter(const QString& outputFilePath);
    ~PlainTextWriter();
    void append(const QString& word, const QString& content);
private:
    QFile* file_;
    QTextStream* out_;
};

#endif // PLAINTEXTWRITER_H
