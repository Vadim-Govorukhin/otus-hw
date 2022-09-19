CREATE TABLE events (
event_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
title VARCHAR(30), 
start_date timestamp with time zone,
end_date timestamp with time zone,
descr TEXT,   
user_id INTEGER,
notify_user_time interval
);