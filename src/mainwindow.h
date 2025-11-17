#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QJsonObject>
#include <QMainWindow>
#include <QNetworkAccessManager>

QT_BEGIN_NAMESPACE
namespace Ui
{
    class MainWindow;
}
class QModelIndex;
QT_END_NAMESPACE

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    MainWindow(QWidget *parent = nullptr);
    ~MainWindow();

private slots:
    void onAddPassword();
    void fetchPasswords();
    void onPasswordItemDoubleClicked(const QModelIndex &index);

private:
    Ui::MainWindow *ui;
    QNetworkAccessManager *m_networkManager;
    QJsonObject m_credentials; // To store the full credential data
};
#endif // MAINWINDOW_H