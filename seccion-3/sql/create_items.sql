CREATE DATABASE labora-proyect-1;

CREATE TABLE public.items
(
    id serial NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    order_date date NOT NULL,
    product VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL CHECK(quantity > 0), 
    price NUMERIC NOT NULL CHECK(price >= 0)
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.items
    OWNER to postgres;