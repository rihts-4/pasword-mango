#include "addeditdialog.h"

#include <QVBoxLayout>
#include <QFormLayout>
#include <QLineEdit>
#include <QPushButton>
#include <QDialogButtonBox>
#include <QNetworkAccessManager>
#include <QNetworkRequest>
#include <QNetworkReply>
#include <QJsonDocument>
#include <QJsonObject>
#include <QMessageBox>
#include <QUrl>

// A basic Ui::AddEditDialog implementation to avoid needing a .ui file
namespace Ui
{
    class AddEditDialog
    {
    public:
        QVBoxLayout *verticalLayout;
        QFormLayout *formLayout;
        QLineEdit *websiteEdit;
        QLineEdit *usernameEdit;
        QLineEdit *passwordEdit;
        QDialogButtonBox *buttonBox;

        void setupUi(QDialog *AddEditDialog)
        {
            AddEditDialog->setObjectName(QString::fromUtf8("AddEditDialog"));
            AddEditDialog->setWindowTitle("Add/Edit Password");
            AddEditDialog->resize(400, 200);
            verticalLayout = new QVBoxLayout(AddEditDialog);
            formLayout = new QFormLayout();
            websiteEdit = new QLineEdit(AddEditDialog);
            usernameEdit = new QLineEdit(AddEditDialog);
            passwordEdit = new QLineEdit(AddEditDialog);
            passwordEdit->setEchoMode(QLineEdit::Password);

            formLayout->addRow("Website:", websiteEdit);
            formLayout->addRow("Username:", usernameEdit);
            formLayout->addRow("Password:", passwordEdit);
            verticalLayout->addLayout(formLayout);

            buttonBox = new QDialogButtonBox(QDialogButtonBox::Ok | QDialogButtonBox::Cancel, AddEditDialog);
            verticalLayout->addWidget(buttonBox);

            // We connect to our own slot to handle submission logic
            // QObject::connect(buttonBox, &QDialogButtonBox::accepted, AddEditDialog, &QDialog::accept);
            QObject::connect(buttonBox, &QDialogButtonBox::rejected, AddEditDialog, &QDialog::reject);
        }
    };
}

AddEditDialog::AddEditDialog(QWidget *parent, const QString &site) : QDialog(parent),
                                                                     ui(new Ui::AddEditDialog),
                                                                     m_networkManager(new QNetworkAccessManager(this)),
                                                                     m_site(site)
{
    ui->setupUi(this);
    connect(ui->buttonBox, &QDialogButtonBox::accepted, this, &AddEditDialog::onAccepted);

    if (!m_site.isEmpty())
    {
        ui->websiteEdit->setText(m_site);
        ui->websiteEdit->setReadOnly(true);
        setWindowTitle("Edit Password");
    }
}

AddEditDialog::~AddEditDialog()
{
    delete ui;
}

QString AddEditDialog::getWebsite() const
{
    return ui->websiteEdit->text().trimmed();
}

QString AddEditDialog::getUsername() const
{
    return ui->usernameEdit->text().trimmed();
}

QString AddEditDialog::getPassword() const
{
    return ui->passwordEdit->text();
}

void AddEditDialog::onAccepted()
{
    QString site = getWebsite();
    QString username = getUsername();
    QString password = getPassword();

    if (site.isEmpty() || username.isEmpty() || password.isEmpty())
    {
        QMessageBox::warning(this, "Input Error", "All fields are required.");
        return;
    }

    QJsonObject json;
    json["username"] = username;
    json["password"] = password;

    QNetworkRequest request;
    QNetworkReply *reply = nullptr;

    if (m_site.isEmpty())
    { // Add new
        request.setUrl(QUrl("http://localhost:8080/credentials"));
        request.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
        json["site"] = site;
        reply = m_networkManager->post(request, QJsonDocument(json).toJson());
    }
    else
    { // Update existing
        request.setUrl(QUrl("http://localhost:8080/credentials/" + m_site));
        request.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
        reply = m_networkManager->put(request, QJsonDocument(json).toJson());
    }

    connect(reply, &QNetworkReply::finished, this, [this, reply]()
            {
        if (reply->error() == QNetworkReply::NoError) {
            int statusCode = reply->attribute(QNetworkRequest::HttpStatusCodeAttribute).toInt();
            if (statusCode >= 200 && statusCode < 300) {
                QMessageBox::information(this, "Success", "Credentials saved successfully.");
                accept(); // Close the dialog on success
            } else {
                QString responseBody = reply->readAll();
                QMessageBox::critical(this, "Server Error",
                                      QString("Failed to save credentials. Status: %1\nResponse: %2")
                                          .arg(statusCode)
                                          .arg(responseBody));
            }
        } else {
            QMessageBox::critical(this, "Network Error",
                                  QString("An error occurred: %1").arg(reply->errorString()));
        }
        reply->deleteLater(); });
}