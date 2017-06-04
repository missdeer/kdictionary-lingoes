/*
 *   Read Lingoes Dictionary Files (*.ld2 or *.ldx)
 *   Copyright (C) 2013-2015 by Symeon Huang <hzwhuang@gmail.com>
 *
 *   This program is free software; you can redistribute it and/or modify
 *   it under the terms of the GNU Library General Public License as
 *   published by the Free Software Foundation; either version 3 or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details
 *
 *   You should have received a copy of the GNU Library General Public
 *   License along with this program; if not, write to the
 *   Free Software Foundation, Inc.,
 *   51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
 */
#include <QCoreApplication>
#include <QFileInfo>
#include <QStringList>
#include <QCommandLineParser>
#include <QDebug>
#include "lingoes.h"

int main(int argc, char** argv)
{
    QCoreApplication app(argc, argv);
    app.setApplicationName("LingoesDictionaryConverter");
    app.setApplicationVersion("1.0");

    QCommandLineParser parser;
    parser.setApplicationDescription("Lingoes dictionary file (LD2/LDX) reader/extracter.");
    parser.addHelpOption();
    parser.addVersionOption();

    QCommandLineOption ldxfile("i", "Input Lingoes dictionary file (default: input.ld2).", "input", "input.ld2");
    QCommandLineOption outfile("o", "Output extracted text file (default: output.txt).", "output", "output.txt");
    QCommandLineOption notrim("disable-trim", "Disable HTML tag trimming (default: no).");
    QCommandLineOption format("f", "Output format, can be plaintext, sqlite or sqlcipher (default: sqlcipher).", "format", "sqlcipher");
    QCommandLineOption autoEncodings("auto-encoding", "Detect encodings automatically (default: no).");
    QCommandLineOption compressed("compressed", "Compress output (default: no).");
    QCommandLineOption cipher("c", "Cipher name, only used when output as sqlcipher (default: aes-256-cbc).", "cipher name", "aes-256-cbc");
    QCommandLineOption key("k", "Cipher key, only used when output as sqlcipher.", "key");

    parser.addOption(ldxfile);
    parser.addOption(outfile);
    parser.addOption(notrim);
    parser.addOption(autoEncodings);
    parser.addOption(format);
    parser.addOption(compressed);
    parser.addOption(cipher);
    parser.addOption(key);

    parser.process(app);

    const QString inputFile = parser.value(ldxfile);
    QFileInfo ld2FileInfo(inputFile);
    if (!ld2FileInfo.exists()) {
        qCritical()<<"Error: Input file" << inputFile << "doesn't exist.";
        return 1;
    }

    if (parser.value(format) == "sqlcipher" && (parser.value(cipher).isEmpty() || parser.value(key).isEmpty()))
    {
        qCritical() << "Need cipher name and key for sqlcipher format.";
        return 1;
    }

    QString ld2file = ld2FileInfo.canonicalFilePath();
    Lingoes ldx(ld2file,
                !parser.isSet(notrim),
                parser.isSet(autoEncodings),
                parser.isSet(compressed),
                parser.value(format),
                parser.value(cipher),
                parser.value(key));
    ldx.extractToFile(parser.value(outfile));

    return 0;
}
