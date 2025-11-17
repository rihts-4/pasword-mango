cmake_minimum_required(VERSION 3.16)
project(pasword-mango)

# Set C++ standard
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Enable Qt's automatic tools
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTOUIC ON)
set(CMAKE_AUTORCC ON)

# Find Qt
find_package(Qt6 REQUIRED COMPONENTS Widgets Network)

# If Qt6 not found, try Qt5
if(NOT Qt6_FOUND)
    find_package(Qt5 REQUIRED COMPONENTS Widgets Network)
endif()

# Create executable
add_executable(${PROJECT_NAME}
    main.cpp
    mainwindow.cpp
    addeditdialog.cpp
    passworddetaildialog.cpp
    passworddetaildialog.h
    ../resources.qrc
)

# Link Qt libraries
target_link_libraries(${PROJECT_NAME}
    PRIVATE
    Qt::Widgets Qt::Network
)

# Set executable properties for Windows
if(WIN32)
    set_target_properties(${PROJECT_NAME} PROPERTIES
        WIN32_EXECUTABLE TRUE
    )
endif()