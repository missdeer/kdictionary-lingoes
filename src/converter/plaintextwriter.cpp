/**************************************************************************
**   Author: Fan Yang
**   Email: missdeer@gmail.com
**   License: see the license.txt file
**************************************************************************/
#include "plaintextwriter.h"

PlainTextWriter::PlainTextWriter(const QString &outputFilePath)
{
    file_ = new QFile(outputFilePath);
    Q_ASSERT(file_);
    if (file_->open(QIODevice::WriteOnly|QIODevice::Text))
    {
        out_ = new QTextStream(file_);
    }
}

PlainTextWriter::~PlainTextWriter()
{
    out_->flush();
    file_->close();
    delete out_;
    delete file_;
}

void PlainTextWriter::append(const QString &word, const QString &content)
{
    *out_ << word << " = " << content << endl;
}
