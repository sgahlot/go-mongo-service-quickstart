// EXECUTE then steps given in this script
use fruit;

db.createUser({
  user: 'test-user',
  pwd: 'test-pass',
  roles: [{
    role: 'readWrite',
    db: 'fruit'
  }]
});

db.fruit.insert({
  _id: '1',
  name: 'Banana',
  description: 'Good for health'
});

db.fruit.insert({
  _id: '2',
  name: 'Apple',
  description: 'Keeps the doctor away'
});

db.fruit.insert({
  _id: '3',
  name: 'Blueberry',
  description: 'Antioxidant Superfood'
});
