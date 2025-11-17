#include "mainwindow.h"
#include "addeditdialog.h"
#include "passworddetaildialog.h"

#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QPushButton>
#include <QListView>
#include <QStringListModel>
#include <QNetworkReply>
#include <QJsonDocument>
#include <QJsonArray>
#include <QMessageBox>

// A basic Ui::MainWindow implementation to avoid needing a .ui file
namespace Ui
{
    class MainWindow
    {
    public:
        QWidget *centralwidget;
        QVBoxLayout *verticalLayout;
        QListView *passwordListView;
        QHBoxLayout *horizontalLayout;
        QPushButton *addButton;
        QPushButton *refreshButton;

        void setupUi(QMainWindow *MainWindow)
        {
            MainWindow->setObjectName(QString::fromUtf8("MainWindow"));
            MainWindow->setWindowTitle("Password Mango");
            MainWindow->resize(500, 400);
            centralwidget = new QWidget(MainWindow);
            verticalLayout = new QVBoxLayout(centralwidget);
            passwordListView = new QListView(centralwidget);
            verticalLayout->addWidget(passwordListView);
            horizontalLayout = new QHBoxLayout();
            addButton = new QPushButton("Add Password", centralwidget);
            refreshButton = new QPushButton("Refresh", centralwidget);
            horizontalLayout->addWidget(addButton);
            horizontalLayout->addWidget(refreshButton);
            verticalLayout->addLayout(horizontalLayout);
            MainWindow->setCentralWidget(centralwidget);
        }
    };
}

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent), ui(new Ui::MainWindow), m_networkManager(new QNetworkAccessManager(this))
{
    ui->setupUi(this);

    connect(ui->addButton, &QPushButton::clicked, this, &MainWindow::onAddPassword);
    connect(ui->refreshButton, &QPushButton::clicked, this, &MainWindow::fetchPasswords);
    connect(ui->passwordListView, &QListView::doubleClicked, this, &MainWindow::onPasswordItemDoubleClicked);

    // Disable editing on double-click to prevent renaming items in the list.
    ui->passwordListView->setEditTriggers(QAbstractItemView::NoEditTriggers);

    fetchPasswords();
}

MainWindow::~MainWindow()
{
    delete ui;
}

void MainWindow::onAddPassword()
{
    AddEditDialog dialog(this);
    if (dialog.exec() == QDialog::Accepted)
    {
        // After adding, refresh the list
        fetchPasswords();
    }
}

void MainWindow::onPasswordItemDoubleClicked(const QModelIndex &index)
{
    QString site = index.data(Qt::DisplayRole).toString();
    // Since we no longer cache the full credential data, we can't check m_credentials.
    // We'll just open the dialog and let it handle fetching or failing.
    if (!site.isEmpty())
    {
        PasswordDetailDialog dialog(site, this);
        // Connect the signal to refresh the list if credentials are changed (e.g., updated or deleted).
        connect(&dialog, &PasswordDetailDialog::credentialsChanged, this, &MainWindow::fetchPasswords);
        dialog.exec();
    }
    else
    {
        // This case is unlikely but good to have.
        QMessageBox::warning(this, "Error", "No site selected.");
    }
}

void MainWindow::fetchPasswords()
{
    QNetworkRequest request(QUrl("http://localhost:8080/credentials"));
    QNetworkReply *reply = m_networkManager->get(request);

    connect(reply, &QNetworkReply::finished, this, [this, reply]()
            {
        if (reply->error() == QNetworkReply::NoError) {
            int statusCode = reply->attribute(QNetworkRequest::HttpStatusCodeAttribute).toInt();
            if (statusCode == 200) {
                QByteArray responseData = reply->readAll();
                QJsonDocument doc = QJsonDocument::fromJson(responseData.trimmed());

                if (!doc.isNull() && doc.isArray()) {
                    QJsonArray sitesArray = doc.array();
                    QStringList sites;
                    for (const QJsonValue &value : sitesArray) {
                        sites.append(value.toString());
                    }
                    sites.sort(Qt::CaseInsensitive);

                    auto *model = qobject_cast<QStringListModel *>(ui->passwordListView->model());
                    if (!model) {
                        model = new QStringListModel(this);
                        ui->passwordListView->setModel(model);
                    }
                    model->setStringList(sites);
                    m_credentials.empty(); // Clear the old cache as it's no longer populated here.
                } else {
                    QMessageBox::warning(this, "Parse Error", "Failed to parse JSON response or it was not a JSON array.");
                }
            } else {
                QMessageBox::warning(this, "Server Error", QString("Failed to fetch passwords. Status: %1").arg(statusCode));
            }
        } else {
            QMessageBox::critical(this, "Network Error", QString("An error occurred: %1").arg(reply->errorString()));
        }
        reply->deleteLater(); });
}