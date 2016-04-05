'use strict';

/**
 * @ngdoc function
 * @name openBicycleDatabaseApp.controller:BikecontrollerCtrl
 * @description
 * # BikecontrollerCtrl
 * Controller of the openBicycleDatabaseApp
 */
angular.module('openBicycleDatabaseApp')
  .controller('BikeListCtrl', function ($scope, Bike) {
    $scope.bikes = Bike.query();
    $scope.remove = function(index,bike){
      console.log(index);
      console.log(bike);
      bike.$remove(function(){
        $scope.bikes = Bike.query();
      });
      
    };
  })
  .controller('BikeDetailCtrl', function ($scope, $routeParams, Bike) {
    Bike.get({id:$routeParams.bikeID},function(bike){
      console.log(bike);
      $scope.bike = bike;
    });
    $scope.bikeSave = function(){
      console.log($scope.bike);
      return $scope.bike.$save();
    };
    $scope.nameLabel = "Name";
    $scope.brandLabel = "Brand";
    $scope.yearLabel = "Year";
    $scope.descriptionLabel = "Description";
    $scope.componentNameLabel = "Component Name";
  })
  .controller('NewBikeCtrl', function($scope, Bike){
    //Label Def
    $scope.nameLabel  = 'Name';
    $scope.brandLabel = 'Brand';
    $scope.yearLabel  = 'Year';
    $scope.namePlaceHolder = 'Bike Name';
    // Errors Def
    $scope.errorBrandRequired = 'The Brand Name is Required';
    $scope.errorYearRequired = 'The Year of the model is Required';
    $scope.errorNameRequired = 'The Name Of the Model is Required';
    $scope.errorYearPattern = ' The Year must be a Digit';
    // Misc
    $scope.validate = 'Submit';

    $scope.save = function(bike){
      console.log(bike);
      Bike.save(bike);
    };
  });
