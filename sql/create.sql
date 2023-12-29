CREATE TABLE Discount (
                          Discount_ID SERIAL PRIMARY KEY,
                          Discount INTEGER,
                          Min_Evaluation DECIMAL(2, 1) CHECK (Min_Evaluation>=0 AND Min_Evaluation<=5)
);

CREATE TABLE "Users" (
                         User_ID SERIAL PRIMARY KEY,
                         Full_Name VARCHAR(50),
                         Email TEXT DEFAULT '',
                         "Password" TEXT ,
                         Phone_Number VARCHAR(20) DEFAULT '',
                         Is_Admin BOOLEAN DEFAULT false,
    --  House_ID INTEGER REFERENCES House(House_ID),
                         UNIQUE (Full_Name)
);

CREATE TABLE Company (
                         Company_ID SERIAL PRIMARY KEY,
                         Name VARCHAR(100),
                         Phone_Number VARCHAR(20),
                         User_ID INTEGER REFERENCES "Users"(User_ID)
);

CREATE TABLE House (
                       House_ID SERIAL PRIMARY KEY,
                       Street TEXT,
                       Number_of_House TEXT,
                       Year_of_Construction INTEGER,
                       Number_of_Floors INTEGER,
                       Number_of_Apartments INTEGER,
                       Evaluation DECIMAL(2, 1),
                       Company_ID INTEGER REFERENCES Company(Company_ID),
                       Discount_ID INTEGER REFERENCES Discount(Discount_ID)

);


CREATE TABLE Complaint (
                           Complaint_ID SERIAL PRIMARY KEY,
                           Title TEXT,
                           Complaint_Text TEXT,
                           Status VARCHAR(20),
                           Evaluation DECIMAL(2, 1) CHECK (Evaluation>=0 AND Evaluation<=5),
                           User_ID INTEGER REFERENCES "Users"(User_ID),
                           Adding_Date DATE
);

CREATE TABLE House_to_Complaint(
                                   Complaint_ID INTEGER REFERENCES Complaint(Complaint_ID),
                                   House_ID INTEGER REFERENCES House(House_ID)
);

CREATE TABLE Company_to_Complaint(
                                     Adding_Evaluation_Date TIMESTAMP,
                                     User_ID INTEGER REFERENCES "Users"(User_ID),
                                     Evaluation_of_Complaint DECIMAL(2, 1),
                                     Complaint_ID INTEGER REFERENCES Complaint(Complaint_ID),
                                     Company_ID INTEGER REFERENCES Company(Company_ID)
);

CREATE TABLE Feedbacks (
                           Feedback_ID SERIAL PRIMARY KEY,
                           User_ID INTEGER REFERENCES "Users"(User_ID),
                           House_ID INTEGER REFERENCES House(House_ID),
                           Review TEXT,
                           Date_adding TIMESTAMP
);

ALTER TABLE "Users" ADD COLUMN House_ID INTEGER REFERENCES House(House_ID);