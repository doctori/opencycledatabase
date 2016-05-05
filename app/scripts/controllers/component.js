'use strict';

/**
 * @ngdoc function
 * @name openBicycleDatabaseApp.controller:ComponentcontrollerCtrl
 * @description
 * # ComponentcontrollerCtrl
 * Controller of the openBicycleDatabaseApp
 */
angular.module('openBicycleDatabaseApp')
  .controller('ComponentListCtrl', function ($scope,Component) {
    $scope.components = Component.query({page:0,per_page:100});
  });
 