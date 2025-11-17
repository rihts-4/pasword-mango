#include <QApplication>
#include <QFile>
#include "mainwindow.h"

int main(int argc, char *argv[])
{
    QApplication app(argc, argv);

    // Load stylesheet
    QFile styleFile(":/styles/app.qss");
    if (styleFile.open(QFile::ReadOnly | QFile::Text))
    {
        QString styleSheet = QLatin1String(styleFile.readAll());
        app.setStyleSheet(styleSheet);
        styleFile.close();
    }

    MainWindow w;
    w.show();

    return app.exec();
}