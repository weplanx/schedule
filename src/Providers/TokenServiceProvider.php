<?php

namespace Lumen\Extra\Providers;

use Laravel\Lumen\Application;
use Illuminate\Support\ServiceProvider;
use Lumen\Extra\Common\TokenFactory;

final class TokenServiceProvider extends ServiceProvider
{
    /**
     * Register the service provider.
     * @return void
     */
    public function register()
    {
        $this->app->singleton('token', function (Application $app) {
            $config = $app->make('config');
            return new TokenFactory(
                $config->get('app.key'),
                $config->get('token')
            );
        });
    }
}
