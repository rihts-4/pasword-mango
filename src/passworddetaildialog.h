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
    void onTogglePasswordVisibility();

private:
    QString m_site;
    QString m_password; // To store the actual password
    bool m_isPasswordVisible;
    QNetworkAccessManager *m_networkManager;
    QLabel *m_usernameLabel;
    QLabel *m_passwordLabel;
    QPushButton *m_togglePasswordButton;
};

#endif // PASSWORDDETAILDIALOG_H