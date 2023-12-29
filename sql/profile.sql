CREATE OR REPLACE FUNCTION get_user_profile(user_id_arg INTEGER)
    RETURNS text
    LANGUAGE SQL AS
$$
SELECT json_build_object(
               'fullname', "Users".Full_Name,
               'phone', "Users".Phone_Number,
               'email', "Users".Email,
               'street', House.Street,
               'number', House.Number_of_House,
               'discount', Discount.Discount
           )
FROM "Users"
         JOIN House ON House.House_ID = "Users".House_ID
         JOIN Discount ON Discount.Discount_ID = House.Discount_ID
WHERE "Users".User_ID = user_id_arg;
$$;

CREATE OR REPLACE FUNCTION update_user_profile(user_id_arg INTEGER, new_fullname TEXT, new_email TEXT, new_phone TEXT, new_street TEXT, new_number_house TEXT)
    RETURNS VOID
    LANGUAGE plpgsql AS
$$
DECLARE
    check_street TEXT;
    check_house TEXT;
    new_house_id INTEGER;
BEGIN
    SELECT House_ID, Street, Number_of_House
    FROM House
    WHERE Street = new_street AND Number_of_House = new_number_house
    INTO new_house_id, check_street, check_house;

    IF new_house_id IS NULL THEN
        RAISE EXCEPTION 'Such house not exist';
    END IF;

    UPDATE "Users"
    SET Full_Name = new_fullname,
        Email = new_email,
        Phone_Number = new_phone,
        House_ID = new_house_id
    WHERE User_ID = user_id_arg;

END;
$$;


CREATE OR REPLACE FUNCTION get_user_reports(p_user_id INTEGER)
    RETURNS text
    LANGUAGE SQL AS
$$
SELECT COALESCE(json_agg(sq), '[]')
FROM (SELECT
          Complaint.Title                          as title,
          Complaint.Complaint_Text                 as complaint_text,
          Complaint.Status                         as complaint_status,
          Complaint.Adding_Date                    as complaint_data

      FROM  "Users"
                JOIN Complaint ON Complaint.User_ID  = "Users".User_ID
      WHERE "Users".User_ID = p_user_id) sq;
$$;

CREATE OR REPLACE FUNCTION add_complaint(p_user_id INTEGER, p_title TEXT, p_complaint_text TEXT)
    RETURNS text
    LANGUAGE plpgsql AS
$$

DECLARE
    p_status TEXT;
    u_id TEXT;
BEGIN
    SELECT User_ID FROM "Users" WHERE User_ID = p_user_id INTO u_id;
    IF u_id IS NULL THEN
        RETURN '';
    END IF;

    INSERT INTO Complaint(Title, Complaint_Text, Status, User_ID, Adding_Date)
    VALUES (p_title, p_complaint_text, p_status, p_user_id, now());
    RETURN '';
END;
$$;

CREATE OR REPLACE FUNCTION set_discount_status()
    RETURNS TRIGGER
    LANGUAGE plpgsql
AS $$
BEGIN
    NEW.Status = 'в обработке';
    RETURN NEW;
END;
$$;

CREATE OR REPLACE TRIGGER set_discount_status_trigger
    BEFORE INSERT
    ON Complaint
    FOR EACH ROW
EXECUTE FUNCTION set_discount_status();

CREATE OR REPLACE FUNCTION add_feedback(p_user_id INTEGER, review_text TEXT)
    RETURNS text
    LANGUAGE plpgsql AS
$$
DECLARE
    p_house_id INTEGER;
    u_id TEXT;
BEGIN
    SELECT User_ID FROM "Users" WHERE User_ID = p_user_id INTO u_id;
    IF u_id IS NULL THEN
        RETURN '';
    END IF;

    SELECT House.House_ID
    FROM "Users"
             JOIN House ON House.House_ID = "Users".House_ID
    WHERE "Users".User_ID = p_user_id
    INTO p_house_id;

    INSERT INTO Feedbacks(User_ID, House_ID, Review, Date_Adding)
    VALUES (p_user_id, p_house_id, review_text, CURRENT_TIMESTAMP);
    RETURN '';
END;
$$;

CREATE OR REPLACE FUNCTION set_user_of_company()
    RETURNS TRIGGER
    LANGUAGE plpgsql AS

$$
BEGIN
    IF new.Is_Admin = true THEN
        UPDATE Company
        SET
            Phone_Number = new.Phone_Number,
            User_ID = new.User_ID
        WHERE Company_ID = (SELECT Company_ID
                            FROM Company WHERE User_ID IS NULL
                            ORDER BY RANDOM()
                            LIMIT 1);
    END IF;
    RETURN NEW;
END;
$$;


CREATE TRIGGER set_user_of_company_trigger
    AFTER INSERT
    ON "Users"
    FOR EACH ROW
EXECUTE FUNCTION set_user_of_company();


CREATE OR REPLACE FUNCTION get_next_report_card(p_user_id INTEGER)
    RETURNS text
    LANGUAGE plpgsql AS
$$
DECLARE
    permission BOOLEAN;
    p_complaint_id INTEGER;
    result TEXT;
BEGIN
    SELECT Is_Admin
    FROM "Users"
    WHERE "Users".User_ID = p_user_id
    INTO permission;


    IF permission = false THEN
        RAISE EXCEPTION 'YOU HAVE NOT PERMISSION';
    END IF;

    SELECT row_to_json(sq)
    FROM(SELECT
             "Users".Full_Name          as full_name,
             Complaint.Adding_Date      as date,
             Complaint.Complaint_ID     as complaint_id,
             Complaint.Title            as title,
             Complaint.Complaint_Text   as text,
             House.Street               as street,
             House.Number_of_House      as number_of_house,
             House.House_ID             as house_id,
             House.Year_of_Construction as year_construct,
             House.Evaluation           as evaluation
         FROM  "Users"
                   JOIN House ON House.House_ID = "Users".House_ID
                   JOIN Complaint ON Complaint.User_ID = "Users".User_ID
                   JOIN Company ON Company.Company_ID = House.Company_ID
         WHERE Company.User_ID = p_user_id) sq
    INTO result;
    RETURN result;
END;
$$;

CREATE OR REPLACE FUNCTION close_complaint(p_complaint_id INTEGER, p_status TEXT, p_evaluation FLOAT)
    RETURNS text
    LANGUAGE plpgsql AS
$$
BEGIN
    UPDATE Complaint
    SET Status = p_status, Evaluation = p_evaluation
    WHERE Complaint.Complaint_ID = p_complaint_id;

    RETURN '';
END;
$$;


CREATE OR REPLACE FUNCTION get_user_company(p_user_id INTEGER)
    RETURNS text
    LANGUAGE plpgsql AS
$$
DECLARE
    result TEXT;
BEGIN
    SELECT row_to_json(sq)
    FROM(SELECT
             Company.Name               as company_name,
             Company.Phone_Number       as phone
         FROM  Company
                   JOIN House ON House.Company_ID = Company.Company_ID
                   JOIN "Users" ON "Users".House_ID = House.House_ID
         WHERE "Users".User_ID = p_user_id) sq
    INTO result;
    RETURN result;
END;
$$;

CREATE OR REPLACE FUNCTION get_is_admin(p_user_id INTEGER)
    RETURNS text
    LANGUAGE plpgsql AS
$$
DECLARE
    permission BOOLEAN;
BEGIN
    SELECT "Users".Is_Admin
    FROM "Users"
    WHERE "Users".User_ID = p_user_id
    INTO permission;
    IF permission = true THEN
        RETURN 'TRUE';
    ELSE
        RETURN 'FALSE';
    END IF;
END;
$$;