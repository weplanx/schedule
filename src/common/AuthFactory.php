<?php

namespace lumen\extra\common;

use lumen\extra\facade\Cookie;
use lumen\extra\facade\Jwt;
use lumen\extra\Redis\RefreshToken;

final class AuthFactory
{
    /**
     * JWT Config
     * @var array $config
     */
    private $config;

    /**
     * AuthFactory constructor.
     * @param array $config
     */
    public function __construct(array $config)
    {
        $this->config = $config;
    }

    /**
     * Get Symbol
     * @param string $scene Scene
     * @return mixed
     * @throws \Exception
     */
    public function symbol(string $scene)
    {
        if (empty($this->config[$scene])) {
            throw new \Exception('not exists scene: ' . $scene);
        }

        if (empty($this->config[$scene]['auth'])) {
            throw new \Exception('must set auth token name');
        }

        $token = Cookie::get($this->config[$scene]['auth']);
        return Jwt::getToken($token)->getClaim('symbol');
    }

    /**
     * Set Cookie Auth Token
     * @param string $scene Scene
     * @param array $symbol Symbol
     * @return bool
     * @throws \Exception
     */
    public function set(string $scene,
                        array $symbol)
    {
        if (empty($this->config[$scene])) {
            throw new \Exception('not exists scene: ' . $scene);
        }

        if (empty($this->config[$scene]['auth'])) {
            throw new \Exception('must set auth token name');
        }

        $token = Jwt::setToken($scene, $symbol);
        if (!$token) {
            return false;
        }

        Cookie::set($this->config[$scene]['auth'], $token);
        return true;
    }

    /**
     * Verify Cookie Auth Token
     * @param string $scene Scene
     * @return bool|string
     * @throws \Exception
     */
    public function verify(string $scene)
    {
        if (empty($this->config[$scene])) {
            throw new \Exception('not exists scene: ' . $scene);
        }

        if (empty($this->config[$scene]['auth'])) {
            throw new \Exception('must set auth token name');
        }


        if (empty(Cookie::get($this->config[$scene]['auth']))) {
            return false;
        }

        $result = Jwt::verify(
            $scene,
            Cookie::get($this->config[$scene]['auth'])
        );

        if (is_string($result)) {
            Cookie::set($this->config[$scene]['auth'], $result);
            return true;
        }

        return $result;
    }

    /**
     * Clear Cookie Auth Token
     * @param string $scene Scene
     * @throws \Exception
     */
    public function clear(string $scene)
    {
        if (empty($this->config[$scene])) {
            throw new \Exception('not exists scene: ' . $scene);
        }

        if (empty($this->config[$scene]['auth'])) {
            throw new \Exception('must set auth token name');
        }

        if (!empty($this->config[$scene]['auto_refresh'])) {
            $token = Jwt::getToken(Cookie::get($this->config[$scene]['auth']));

            $result = (new RefreshToken)->clear(
                $token->getClaim('jti'),
                $token->getClaim('ack')
            );

            if (!$result) {
                throw new \Exception('clear refresh token failed');
            }
        }

        Cookie::delete($this->config[$scene]['auth']);
    }

}
