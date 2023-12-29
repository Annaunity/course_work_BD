CREATE OR REPLACE FUNCTION get_house_stat(p_user_id INTEGER)
    RETURNS text
    LANGUAGE SQL AS
$$
SELECT row_to_json(sq)
FROM (SELECT
          House.Street                             as street_house,
          House.Number_of_House                    as number_house,
          House.Year_of_Construction               as year_construct,
          House.Number_of_Floors                   as num_of_floors,
          House.Number_of_Apartments               as number_of_apartments,
          House.Evaluation                         as evaluation,
          Discount.Discount                        as discount
      FROM  "Users"
                JOIN House ON House.House_ID = "Users".House_ID
                JOIN Discount ON House.Discount_ID = Discount.Discount
      WHERE "Users".User_ID = p_user_id) sq;

$$;

CREATE OR REPLACE FUNCTION get_house_feedbacks(p_user_id INTEGER)
    RETURNS text
    LANGUAGE SQL AS
$$
SELECT COALESCE(json_agg(sq), '[]')
FROM (SELECT
          Feedbacks.Review                         as feedback_text,
          Feedbacks.Date_adding                    as feedback_date

      FROM  "Users"
                JOIN House ON House.House_ID = "Users".House_ID
                JOIN Feedbacks ON Feedbacks.House_ID  = House.House_ID
      WHERE "Users".User_ID = p_user_id) sq;
$$;

CREATE OR REPLACE FUNCTION add_house(p_street TEXT, p_number_house TEXT, p_year_of_construction INTEGER, p_number_of_floors INTEGER, p_number_of_apartments INTEGER)
    RETURNS INTEGER
    LANGUAGE plpgsql AS
$$

DECLARE
    new_house_id INTEGER := 0;
    first_evaluation DECIMAL(2, 1) := 5.0;
    p_company_id INTEGER;
BEGIN
    SELECT Company_ID
    FROM Company
    ORDER BY RANDOM()
    LIMIT 1
    INTO p_company_id;

    INSERT INTO House(Street, Number_of_House, Year_of_Construction, Number_of_Floors, Number_of_Apartments, Evaluation, Company_ID)
    VALUES (p_street, p_number_house, p_year_of_construction, p_number_of_floors, p_number_of_apartments, first_evaluation, p_company_id)
    RETURNING House_ID INTO new_house_id;
    RETURN new_house_id;
END;
$$;

CREATE OR REPLACE FUNCTION calculate_discount()
    RETURNS TRIGGER
    LANGUAGE plpgsql
AS $$
BEGIN
    NEW.Discount_ID = (SELECT Discount.Discount_ID
                       FROM Discount
                       WHERE NEW.Evaluation >= Min_Evaluation
                       ORDER BY -Discount.Discount
                       LIMIT 1);
    RETURN NEW;
END;
$$;

CREATE OR REPLACE TRIGGER discount_to_house
    BEFORE INSERT OR UPDATE
    ON House
    FOR EACH ROW
EXECUTE PROCEDURE calculate_discount();


CREATE OR REPLACE FUNCTION new_evaluation() RETURNS TRIGGER AS
$$
DECLARE
    avg_evaluation INTEGER;
    p_house_id INTEGER;
BEGIN
    SELECT AVG(Complaint.Evaluation)
    FROM Complaint
             JOIN "Users" ON Complaint.User_ID = "Users".User_ID
             JOIN House ON House.House_ID = "Users".House_ID
    WHERE Complaint.Status <> 'в обработке'
      AND "Users".User_ID = Complaint.User_ID
      AND House.House_ID = "Users".House_ID
    INTO avg_evaluation;

    SELECT House.house_id FROM Complaint
        JOIN "Users" ON Complaint.user_id = "Users".user_id
        JOIN  House ON House.house_id = "Users".house_id
    WHERE complaint_id = new.complaint_id
    INTO p_house_id;

    UPDATE House
    SET Evaluation = avg_evaluation
    WHERE  House_ID = p_house_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE TRIGGER complaint_to_evaluation
    AFTER UPDATE
    ON Complaint
    FOR EACH ROW
EXECUTE PROCEDURE new_evaluation();