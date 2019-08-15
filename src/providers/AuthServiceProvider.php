<?php

namespace lumen\extra\providers;

use Laravel\Lumen\Application;
use Illuminate\Support\ServiceProvider;
use lumen\extra\common\AuthFactory;
use lumen\extra\common\JwtFactory;

class AuthServiceProvider extends ServiceProvider
{
    /**
     * Register the service provider.
     *
     * @return void
     */
    public function register()
    {
        $this->app->singleton('auth', function (Application $app) {
            $config = $app->make('config')
                ->get('jwt');

            return new AuthFactory($config);
        });
    }
}
