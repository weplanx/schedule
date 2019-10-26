<?php

namespace lumen\extra\providers;

use Laravel\Lumen\Application;
use Illuminate\Support\ServiceProvider;
use lumen\extra\common\TokenFactory;

class TokenServiceProvider extends ServiceProvider
{
    /**
     * Register the service provider.
     *
     * @return void
     */
    public function register()
    {
        $this->app->singleton('token', function (Application $app) {
            $secret = $app->make('config')
                ->get('app.key');
            $config = $app->make('config')
                ->get('token');

            return new TokenFactory($secret, $config);
        });
    }
}
