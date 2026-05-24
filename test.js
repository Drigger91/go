var num = 10; // Global
const myObject = { num: 50 };

function test() {
  console.log(this.num);
}

test(); 
test.call(myObject);