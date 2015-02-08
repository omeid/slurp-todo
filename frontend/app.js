var app = angular.module('app', ['ngResource']);

app.factory('TODO', function($resource) {
  return $resource('api/todos/:todoId', {
    todoId: '@id'
  }, {
    update: { method: 'PUT' }
  });
});


app.controller('todos', function($scope, TODO) {
  
  $scope.todos = TODO.query()
  
  $scope.getTotalTodos = function () {
    return $scope.todos.length;
  };
  
  
  $scope.addTodo = function () {
    var todo = new TODO({task:$scope.formTodoText, done:false});
    todo.$save();
    $scope.todos.push(todo);

    $scope.formTodoText = '';
  };

  $scope.removeTodo = function(index) {
      $scope.todos[index].$delete()
      $scope.todos.splice(index, 1)
  }
  
  $scope.updateTodo = function(index) {
      $scope.todos[index].$update()
  }

  $scope.clearCompleted = function () {
    $scope.todos = _.filter($scope.todos, function(todo){
      return !todo.done;
    });
  };
})

// Small helper functions to include javascript.
function require(path) {
  var js = document.createElement("script");
  js.type = "text/javascript";
  js.src = path;
  document.body.appendChild(js);
}
