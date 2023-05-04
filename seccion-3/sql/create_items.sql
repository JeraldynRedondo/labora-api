CREATE DATABASE labora-proyect-1;

CREATE TABLE public.items
(
    id serial NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    order_date date NOT NULL,
    product VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    price NUMERIC NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.items
    OWNER to postgres;