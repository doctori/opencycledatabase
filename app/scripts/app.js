'use strict';

/**
 * @ngdoc overview
 * @name openBicycleDatabaseApp
 * @description
 * # openBicycleDatabaseApp
 *
 * Main module of the application.
 */
angular
  .module('openBicycleDatabaseApp', [
    'bikeService',
    'ngAnimate',
    'ngCookies',
    'ngResource',
    'ngRoute',
    'ngSanitize',
    'ngTouch'
  ])
  .config(function ($routeProvider) {
    $routeProvider
      .when('/', {
        templateUrl: 'views/bike.html',
        controller: 'BikeListCtrl',
        controllerAs: 'bike'
      })
      .when('/about', {
        templateUrl: 'views/about.html',
        controller: 'AboutCtrl',
        controllerAs: 'about'
      })
      .otherwise({
        redirectTo: '/'
      });
  });