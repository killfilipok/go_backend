<html>
    <head>
    <title></title>
    </head>
    <body>
        <form action="/login" method="post">
            Username:<input type="text" name="username">
            Password:<input type="password" name="password">
            <input type="submit" value="Login">
            <input type="radio" name="gender" value="1">Male
    <input type="radio" name="gender" value="2">Female
            <select name="fruit">
    <option value="apple">apple</option>
    <option value="pear">pear</option>
    <option value="banana">banana</option>
    </select>
        <input type="hidden" name="token" value="{{.}}">
    <input type="checkbox" name="interest" value="football">Football
    <input type="checkbox" name="interest" value="basketball">Basketball
    <input type="checkbox" name="interest" value="tennis">Tennis
        </form>

    </body>
</html>