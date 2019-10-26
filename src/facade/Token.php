<?php

namespace lumen\extra\facade;

use Lcobucci\JWT\Token as LcobucciToken;
use Illuminate\Support\Facades\Facade;

/**
 * Class Token
 * @method static LcobucciToken getToken(string $token = null)
 * @method static bool|string setToken(string $scene, array $symbol = [])
 * @method static bool|string verify(string $scene, string $token)
 * @package lumen\extra\facade
 */
class Token extends Facade
{
    /**
     * Get the registered name of the component.
     *
     * @return string
     */
    protected static function getFacadeAccessor()
    {
        return 'token';
    }
}
