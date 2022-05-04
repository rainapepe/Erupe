1. Extract into Erupe, overwrite anything that exists.
1a (Optional). Change subscription type in server/channelserver/handlers.go
2. Open pgadmin
3. Erupe > Public > Tables > Right-click on Users > Query Tool
4. Paste the following text in there.

ALTER TABLE IF EXISTS public.users
    ADD COLUMN item_box bytea;

5. Press F5.
6. If it says "Query successful" then you can just close pgadmin.
7. Run the server and enjoy the Delivery Cat.