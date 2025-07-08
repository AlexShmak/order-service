revoke select, insert, update, delete on orders_service.users from admin_user;
revoke select on orders_service.users from regular_user;
revoke usage on schema orders_service from regular_user, admin_user;
drop role if exists admin_user, regular_user;
