<?php

namespace lumen\extra\facade;

use Lcobucci\JWT\Token;
use Illuminate\Support\Facades\Facade;

/**
 * Class Jwt
 * @method static Token getToken()
 * @method static bool|string setToken(string $scene, array $symbol = [])
 * @method static bool|string verify(string $scene, string $token)
 * @package lumen\extra\jwt
 */
class Jwt extends Facade
{
    /**
     * Get the registered name of the component.
     *
     * @return string
     */
    protected static function getFacadeAccessor()
    {
        return 'jwt';
    }
}
