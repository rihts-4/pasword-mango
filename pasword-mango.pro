QT += core gui widgets network

CONFIG += c++17

TARGET = pasword-mango
TEMPLATE = app

SOURCES += \
    src/main.cpp \
    src/mainwindow.cpp \
    src/addeditdialog.cpp

HEADERS += \
    src/mainwindow.h \
    src/addeditdialog.h

RESOURCES += \
    resources.qrc

INCLUDEPATH += \
    src