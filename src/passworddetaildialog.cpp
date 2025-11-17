#include "passworddetaildialog.h"
#include "addeditdialog.h"

#include <QVBoxLayout>
#include <QFormLayout>
#include <QLabel>
#include <QPushButton>
#include <QDialogButtonBox>
#include <QMessageBox>
#include <QNetworkAccessManager>
#include <QNetworkRequest>
#include <QNetworkReply>
#include <QJsonDocument>
#include <QJsonObject>
#include <QUrl>
#include <QDebug>

PasswordDetailDialog::PasswordDetailDialog(const QString &site, QWidget *parent)
    : QDialog(parent),
      m_site(site),
      m_isPasswordVisible(false),
      m_networkManager(new QNetworkAccessManager(this))
{
    setWindowTitle("Password Details");
    QVBoxLayout *layout = new QVBoxLayout(this);
    QFormLayout *formLayout = new QFormLayout();

    m_usernameLabel = new QLabel("Loading...", this);
    m_passwordLabel = new QLabel("Loading...", this);
    m_passwordLabel->setTextInteractionFlags(Qt::TextSelectableByMouse);
    m_usernameLabel->setTextInteractionFlags(Qt::TextSelectableByMouse);
    formLayout->addRow("Site:", new QLabel(m_site, this));
    formLayout->addRow("Username:", m_usernameLabel);

    // Layout for password and toggle button
    QHBoxLayout *passwordLayout = new QHBoxLayout();
    passwordLayout->addWidget(m_passwordLabel);
    m_togglePasswordButton = new QPushButton("Show", this);
    m_togglePasswordButton->setFixedWidth(60);
    passwordLayout->addWidget(m_togglePasswordButton);
    formLayout->addRow("Password:", passwordLayout);

    layout->addLayout(formLayout);

    QDialogButtonBox *buttonBox = new QDialogButtonBox(this);
    QPushButton *updateButton = buttonBox->addButton("Update", QDialogButtonBox::ActionRole);
    QPushButton *deleteButton = buttonBox->addButton("Delete", QDialogButtonBox::DestructiveRole);
    buttonBox->addButton(QDialogButtonBox::Close);

    layout->addWidget(buttonBox);

    connect(updateButton, &QPushButton::clicked, this, &PasswordDetailDialog::onUpdate);
    connect(deleteButton, &QPushButton::clicked, this, &PasswordDetailDialog::onDelete);
    connect(buttonBox, &QDialogButtonBox::rejected, this, &QDialog::reject);
    connect(m_togglePasswordButton, &QPushButton::clicked, this, &PasswordDetailDialog::onTogglePasswordVisibility);

    // Fetch the specific credentials for this site
    QNetworkRequest request(QUrl("http://localhost:8080/credentials/" + m_site));
    QNetworkReply *reply = m_networkManager->get(request);
    connect(reply, &QNetworkReply::finished, this, [this, reply]()
            { onCredentialsFetched(reply); });
}

PasswordDetailDialog::~PasswordDetailDialog()
{
    // The QObject parent-child relationship will handle deletion of m_networkManager and labels.
}

void PasswordDetailDialog::onCredentialsFetched(QNetworkReply *reply)
{
    if (reply->error() == QNetworkReply::NoError)
    {
        QJsonDocument doc = QJsonDocument::fromJson(reply->readAll());

        if (doc.isObject())
        {
            QJsonObject creds = doc.object();
            // Use case-insensitive keys from Go backend
            m_password = creds["Password"].toString();
            m_usernameLabel->setText(creds["Username"].toString());
            m_passwordLabel->setText("******");
        }
        else
        {
            QMessageBox::critical(this, "Error", "Failed to parse credential details from server response.");
            reject(); // Close dialog on error
        }
    }
    else
    {
        QMessageBox::critical(this, "Network Error", "Failed to fetch credential details: " + reply->errorString());
        reject(); // Close dialog on error
    }
    reply->deleteLater();
}

void PasswordDetailDialog::onTogglePasswordVisibility()
{
    if (!m_isPasswordVisible)
    {
        auto reply = QMessageBox::question(this, "Show Password",
                                           "Are you sure you want to show the password in plain text?",
                                           QMessageBox::Yes | QMessageBox::No);

        if (reply == QMessageBox::Yes)
        {
            m_passwordLabel->setText(m_password);
            m_togglePasswordButton->setText("Hide");
            m_isPasswordVisible = true;
        }
    }
    else
    {
        m_passwordLabel->setText("******");
        m_togglePasswordButton->setText("Show");
        m_isPasswordVisible = false;
    }
}

void PasswordDetailDialog::onUpdate()
{
    AddEditDialog dialog(this, m_site);
    if (dialog.exec() == QDialog::Accepted)
    {
        emit credentialsChanged();
        accept(); // Close this dialog
    }
}

void PasswordDetailDialog::onDelete()
{
    auto reply = QMessageBox::question(this, "Confirm Delete",
                                       QString("Are you sure you want to delete the credentials for %1?").arg(m_site),
                                       QMessageBox::Yes | QMessageBox::No);

    if (reply == QMessageBox::No)
    {
        return;
    }

    QNetworkRequest request(QUrl("http://localhost:8080/credentials/" + m_site));
    QNetworkReply *deleteReply = m_networkManager->deleteResource(request);

    connect(deleteReply, &QNetworkReply::finished, this, [this, deleteReply]()
            {
        if (deleteReply->error() == QNetworkReply::NoError) {
            emit credentialsChanged();
            accept(); // Close the dialog
        } else {
            QMessageBox::critical(this, "Error", "Failed to delete credentials: " + deleteReply->errorString());
        }
        deleteReply->deleteLater(); });
}