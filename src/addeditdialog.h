#ifndef ADDEDITDIALOG_H
#define ADDEDITDIALOG_H

#include <QDialog>
class QNetworkAccessManager;

QT_BEGIN_NAMESPACE
namespace Ui
{
    class AddEditDialog;
}
QT_END_NAMESPACE

class AddEditDialog : public QDialog
{
    Q_OBJECT

public:
    explicit AddEditDialog(QWidget *parent = nullptr, const QString &site = QString());
    ~AddEditDialog();

    QString getWebsite() const;
    QString getUsername() const;
    QString getPassword() const;

private slots:
    void onAccepted();

private:
    Ui::AddEditDialog *ui;
    QNetworkAccessManager *m_networkManager;
    QString m_site; // Used to determine if we are in "edit" mode
};

#endif // ADDEDITDIALOG_H