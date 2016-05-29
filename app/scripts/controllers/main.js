'use strict';

/**
 * @ngdoc function
 * @name openBicycleDatabaseApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the openBicycleDatabaseApp
 */
angular.module('openBicycleDatabaseApp')
  .controller('MainCtrl', function ($scope,$mdSidenav) {
    this.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];
    $scope.toggleBikeList = function () {
      console.log("TOGGLE");
      $mdSidenav('right').toggle();
    };
  });
