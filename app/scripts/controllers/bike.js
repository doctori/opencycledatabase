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
    this.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];
    $scope.bikes = Bike.query(function(){
    	console.log(bikes);
    });

  });
