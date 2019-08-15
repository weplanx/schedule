<?php

namespace lumen\extra\facade;

use Illuminate\Support\Facades\Facade;

/**
 * Class Auth
 * @method static mixed symbol(string $scene)
 * @method static bool set(string $scene, array $symbol)
 * @method static bool|string verify(string $scene)
 * @method static void clear(string $scene)
 * @package lumen\extra\facade
 */
class Auth extends Facade
{
    /**
     * Get the registered name of the component.
     *
     * @return string
     */
    protected static function getFacadeAccessor()
    {
        return 'auth';
    }
}
