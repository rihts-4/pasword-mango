#ifndef PASSWORDDETAILDIALOG_H
#define PASSWORDDETAILDIALOG_H

#include <QDialog>
#include <QJsonObject>

class QNetworkAccessManager;
class QLabel;
class QNetworkReply;

class PasswordDetailDialog : public QDialog
{
    Q_OBJECT

public:
    explicit PasswordDetailDialog(const QString &site, QWidget *parent = nullptr);
    ~PasswordDetailDialog();

signals:
    void credentialsChanged();

private slots:
    void onUpdate();
    void onDelete();
    void onCredentialsFetched(QNetworkReply *reply);

private:
    QString m_site;
    QNetworkAccessManager *m_networkManager;
    QLabel *m_usernameLabel;
    QLabel *m_passwordLabel;
};

#endif // PASSWORDDETAILDIALOG_H