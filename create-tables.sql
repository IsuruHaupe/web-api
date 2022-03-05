use web_api_database;

DROP TABLE IF EXISTS person_has_contact;

DROP TABLE IF EXISTS contact_have_skills;

DROP TABLE IF EXISTS person;

DROP TABLE IF EXISTS skills;

CREATE TABLE person (
  id INT PRIMARY KEY NOT NULL AUTO_INCREMENT, 
  firstname VARCHAR(255),
  lastname VARCHAR(255), 
  fullname VARCHAR(255), 
  home_address VARCHAR(255), 
  email VARCHAR(255), 
  phone_number VARCHAR(255)
);

CREATE TABLE skills (
  id INT PRIMARY KEY NOT NULL AUTO_INCREMENT, 
  skill_name VARCHAR(50), 
  skill_level ENUM('Familiar', 'Proficient', 'Excellent', 'Expert')
);

CREATE TABLE person_has_contact (
  person_id INT NOT NULL, 
  contact_id INT NOT NULL, 
  CONSTRAINT has_contact FOREIGN KEY (person_id) REFERENCES person(id) ON DELETE CASCADE,
  CONSTRAINT is_contact FOREIGN KEY (contact_id) REFERENCES person(id) ON DELETE CASCADE
);

CREATE TABLE contact_have_skills (
  person_id INT NOT NULL, 
  skill_id INT NOT NULL, 
  CONSTRAINT personx FOREIGN KEY (person_id) REFERENCES person(id) ON DELETE CASCADE,
  CONSTRAINT skillx FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);

INSERT INTO person (firstname, lastname, fullname, home_address, email, phone_number)
 VALUES
 ('test1', 'TEST1', 'test1 TEST1', 'Adress1', 'email@email.com', '023209382'),
 ('test2', 'TEST2', 'test2 TEST2', 'Adress2', 'email@email.com', '023209382'),
 ('test3', 'TEST3', 'test3 TEST3', 'Adress3', 'email@email.com', '023209382'),
 ('test4', 'TEST4', 'test4 TEST4', 'Adress4', 'email@email.com', '023209382');


 INSERT INTO skills (skill_name, skill_level) 
 VALUES
   ('Go', 'Proficient'), 
   ('Scala', 'Proficient'), 
   ('Java', 'Excellent');

INSERT INTO person_has_contact (person_id, contact_id)
VALUES
(1, 2),
(2, 3),
(1, 3);

INSERT INTO contact_have_skills (person_id, skill_id)
VALUES
(1, 1),
(1, 2),
(2, 1);