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
    'ngTouch',
    'xeditable'
  ])
  .config(function ($routeProvider) {
    $routeProvider
      .when('/', {
        templateUrl: 'views/bike.html',
        controller: 'BikeListCtrl',
        controllerAs: 'bike'
      })
      .when('/bike/:bikeID',{
        templateUrl: 'views/bike-detail.html',
        controller: 'BikeDetailCtrl',
        controllerAs: 'bike'
      })
      .when('/new-bike',{
        templateUrl: 'views/new-bike.html',
        controller: 'NewBikeCtrl',
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