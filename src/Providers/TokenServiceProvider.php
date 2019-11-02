<?php

namespace Lumen\Extra\Providers;

use Laravel\Lumen\Application;
use Illuminate\Support\ServiceProvider;
use Lumen\Extra\Common\TokenFactory;
use Lumen\Extra\Contracts\TokenInterface;

final class TokenServiceProvider extends ServiceProvider
{
    /**
     * Register the service provider.
     * @return void
     */
    public function register()
    {
        $this->app->singleton(TokenInterface::class, function (Application $app) {
            $config = $app->make('config');
            return new TokenFactory(
                $config->get('app.key'),
                $config->get('token')
            );
        });
    }
}
