create role regular_user;
create role admin_user;

grant usage on schema orders_service to regular_user, admin_user;
grant select, insert, update, delete on orders_service.users to admin_user;
grant select on orders_service.users to regular_user;